package service

import (
	"context"
	"inventory-service/config"
	"inventory-service/internal/adapter/repository"
	postgresrepository "inventory-service/internal/adapter/repository/postgres"
	"inventory-service/internal/domain/entity"
	"inventory-service/pkg/logger"
)

var _ ReservationService = (*reservationService)(nil)

type ReservationService interface {
	Find(ctx context.Context, filter *postgresrepository.FilterReservationPayload) ([]*entity.Reservation, int, error)
	FindByID(ctx context.Context, id uint32) (*entity.Reservation, error)
	Create(ctx context.Context, reservation *entity.Reservation) (*entity.Reservation, error)
	UpdateStatus(ctx context.Context, ids []uint32, status string) error
}

type reservationService struct {
	config *config.Config
	repo   repository.Repository
	logger logger.Logger
}

func NewReservationService(config *config.Config, repo repository.Repository, logger logger.Logger) *reservationService {
	return &reservationService{config: config, repo: repo, logger: logger}
}

func (s *reservationService) Find(ctx context.Context, filter *postgresrepository.FilterReservationPayload) ([]*entity.Reservation, int, error) {
	return s.repo.Postgres().Reservation().Find(ctx, filter)
}

func (s *reservationService) FindByID(ctx context.Context, id uint32) (*entity.Reservation, error) {
	return s.repo.Postgres().Reservation().FindByID(ctx, id)
}

func (s *reservationService) Create(ctx context.Context, reservation *entity.Reservation) (*entity.Reservation, error) {
	var createdReservation *entity.Reservation

	atomic := func(txRepo postgresrepository.PostgresRepository) error {
		var err error
		createdReservation, err = txRepo.Reservation().Create(ctx, reservation)
		return err
	}

	err := s.repo.Postgres().Atomic(ctx, s.config, atomic)
	if err != nil {
		return nil, err
	}

	return createdReservation, nil
}

func (s *reservationService) UpdateStatus(ctx context.Context, ids []uint32, status string) error {
	atomic := func(txRepo postgresrepository.PostgresRepository) error {
		return txRepo.Reservation().UpdateStatus(ctx, ids, status)
	}

	err := s.repo.Postgres().Atomic(ctx, s.config, atomic)
	if err != nil {
		return err
	}

	return nil
}
