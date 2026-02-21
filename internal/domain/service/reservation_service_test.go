package service_test

import (
	"context"
	"testing"

	"inventory-service/config"
	postgresrepository "inventory-service/internal/adapter/repository/postgres"
	"inventory-service/internal/domain/entity"
	"inventory-service/internal/domain/service"
	"inventory-service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper to initialize the mock chain for Reservations
func setupReservationMocks(t *testing.T) (*mocks.MockRepository, *mocks.MockPostgresRepository, *mocks.MockReservationRepository) {
	mRepo := mocks.NewMockRepository(t)
	mPostgres := mocks.NewMockPostgresRepository(t)
	mReservation := mocks.NewMockReservationRepository(t)

	// Link Repository -> PostgresRepository
	mRepo.EXPECT().Postgres().Return(mPostgres).Maybe()
	// Link PostgresRepository -> ReservationRepository
	mPostgres.EXPECT().Reservation().Return(mReservation).Maybe()

	return mRepo, mPostgres, mReservation
}

func TestReservationServiceFind(t *testing.T) {
	mockRepo, _, mockRes := setupReservationMocks(t)
	ctx := context.Background()
	filter := &postgresrepository.FilterReservationPayload{Page: 1}
	expectedList := []*entity.Reservation{&entity.Reservation{Base: entity.Base{ID: 1}}}

	mockRes.EXPECT().Find(ctx, filter).Return(expectedList, 1, nil)

	resService := service.NewReservationService(service.Properties{Repo: mockRepo})
	results, total, err := resService.Find(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, results, 1)
}

func TestReservationServiceFindByID(t *testing.T) {
	mockRepo, _, mockRes := setupReservationMocks(t)
	ctx := context.Background()
	id := uint32(1)
	expected := &entity.Reservation{Base: entity.Base{ID: id}}

	mockRes.EXPECT().FindByID(ctx, id).Return(expected, nil)

	resService := service.NewReservationService(service.Properties{Repo: mockRepo})
	result, err := resService.FindByID(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, id, result.ID)
}

func TestReservationServiceCreateAtomic(t *testing.T) {
	mockRepo, mockPostgres, mockRes := setupReservationMocks(t)
	ctx := context.Background()
	input := &entity.Reservation{ProductID: 10}
	expected := &entity.Reservation{Base: entity.Base{ID: 1}, ProductID: 10}

	// 1. Mock the Atomic call
	// We use Run to execute the callback passed to Atomic
	mockPostgres.EXPECT().
		Atomic(ctx, mock.Anything, mock.Anything).
		Run(func(ctx context.Context, cfg *config.Config, fn postgresrepository.RepositoryAtomicCallback) {
			// Execute the callback using the mockPostgres so internal calls work
			_ = fn(mockPostgres)
		}).
		Return(nil)

	// 2. Mock the Create call inside the atomic block
	mockRes.EXPECT().Create(ctx, input).Return(expected, nil)

	resService := service.NewReservationService(service.Properties{
		Repo:   mockRepo,
		Config: &config.Config{},
	})

	result, err := resService.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint32(1), result.ID)
}

func TestReservationServiceUpdateStatusAtomic(t *testing.T) {
	mockRepo, mockPostgres, mockRes := setupReservationMocks(t)
	ctx := context.Background()
	ids := []uint32{1, 2}
	status := "COMPLETED"

	// Mock Atomic transaction
	mockPostgres.EXPECT().
		Atomic(ctx, mock.Anything, mock.Anything).
		Run(func(ctx context.Context, cfg *config.Config, fn postgresrepository.RepositoryAtomicCallback) {
			_ = fn(mockPostgres)
		}).
		Return(nil)

	// Mock UpdateStatus inside the transaction
	mockRes.EXPECT().UpdateStatus(ctx, ids, status).Return(nil)

	resService := service.NewReservationService(service.Properties{
		Repo:   mockRepo,
		Config: &config.Config{},
	})

	err := resService.UpdateStatus(ctx, ids, status)

	assert.NoError(t, err)
}
