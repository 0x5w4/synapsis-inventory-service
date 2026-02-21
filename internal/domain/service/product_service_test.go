package service_test

import (
	"context"
	"testing"

	postgresrepository "inventory-service/internal/adapter/repository/postgres"
	"inventory-service/internal/domain/entity"
	"inventory-service/internal/domain/service"
	"inventory-service/mocks"

	"github.com/stretchr/testify/assert"
)

// Helper function to initialize the mock chain
func setupProductMocks(t *testing.T) (*mocks.MockRepository, *mocks.MockPostgresRepository, *mocks.MockProductRepository) {
	mRepo := mocks.NewMockRepository(t)
	mPostgres := mocks.NewMockPostgresRepository(t)
	mProduct := mocks.NewMockProductRepository(t)

	// Link the layers together as defined in repository.go and postgres.go
	mRepo.EXPECT().Postgres().Return(mPostgres).Maybe()
	mPostgres.EXPECT().Product().Return(mProduct).Maybe()

	return mRepo, mPostgres, mProduct
}

func TestProductServiceCreate(t *testing.T) {
	mockRepo, _, mockProduct := setupProductMocks(t)

	ctx := context.Background()
	input := &entity.Product{Name: "Test Product"}
	expectedOutput := &entity.Product{ID: 1, Name: "Test Product"}

	// Mock the call on the leaf repository
	mockProduct.EXPECT().Create(ctx, input).Return(expectedOutput, nil)

	productService := service.NewProductService(service.Properties{Repo: mockRepo})
	result, err := productService.Create(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, uint32(1), result.ID)
}

func TestProductServiceUpdate(t *testing.T) {
	mockRepo, _, mockProduct := setupProductMocks(t)

	ctx := context.Background()
	input := &entity.Product{ID: 1, Name: "Updated Product"}

	mockProduct.EXPECT().Update(ctx, input).Return(input, nil)

	productService := service.NewProductService(service.Properties{Repo: mockRepo})
	result, err := productService.Update(ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Product", result.Name)
}

func TestProductServiceDelete(t *testing.T) {
	mockRepo, _, mockProduct := setupProductMocks(t)

	ctx := context.Background()
	id := uint32(1)

	mockProduct.EXPECT().Delete(ctx, id).Return(nil)

	productService := service.NewProductService(service.Properties{Repo: mockRepo})
	err := productService.Delete(ctx, id)

	assert.NoError(t, err)
}

func TestProductServiceFind(t *testing.T) {
	mockRepo, _, mockProduct := setupProductMocks(t)

	ctx := context.Background()
	filter := &postgresrepository.FilterProductPayload{Page: 1, PerPage: 10}
	expectedList := []*entity.Product{{ID: 1, Name: "Item A"}}

	mockProduct.EXPECT().Find(ctx, filter).Return(expectedList, 1, nil)

	productService := service.NewProductService(service.Properties{Repo: mockRepo})
	products, total, err := productService.Find(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, products, 1)
}

func TestProductServiceFindByID(t *testing.T) {
	mockRepo, _, mockProduct := setupProductMocks(t)

	ctx := context.Background()
	id := uint32(1)
	expected := &entity.Product{ID: 1, Name: "Item A"}

	mockProduct.EXPECT().FindByID(ctx, id).Return(expected, nil)

	productService := service.NewProductService(service.Properties{Repo: mockRepo})
	result, err := productService.FindByID(ctx, id)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, id, result.ID)
}
