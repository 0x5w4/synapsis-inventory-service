package service

import (
	"context"
	postgresrepository "inventory-service/internal/adapter/repository/postgres"
	"inventory-service/internal/domain/entity"
)

var _ ReservationService = (*reservationService)(nil)

type ReservationService interface {
	Find(ctx context.Context, filter *postgresrepository.FilterReservationPayload) ([]*entity.Reservation, int, error)
	FindByID(ctx context.Context, id uint32) (*entity.Reservation, error)
	Create(ctx context.Context, reservation *entity.Reservation) (*entity.Reservation, error)
	UpdateStatus(ctx context.Context, ids []uint32, status string) error
}

type reservationService struct {
	Properties
}

func NewReservationService(props Properties) *reservationService {
	return &reservationService{Properties: props}
}

func (s *reservationService) Find(ctx context.Context, filter *postgresrepository.FilterReservationPayload) ([]*entity.Reservation, int, error) {
	return s.Repo.Postgres().Reservation().Find(ctx, filter)
}

func (s *reservationService) FindByID(ctx context.Context, id uint32) (*entity.Reservation, error) {
	return s.Repo.Postgres().Reservation().FindByID(ctx, id)
}

func (s *reservationService) Create(ctx context.Context, reservation *entity.Reservation) (*entity.Reservation, error) {
	var createdReservation *entity.Reservation

	atomic := func(txRepo postgresrepository.PostgresRepository) error {
		var err error
		createdReservation, err = txRepo.Reservation().Create(ctx, reservation)
		return err
	}

	err := s.Repo.Postgres().Atomic(ctx, s.Config, atomic)
	if err != nil {
		return nil, err
	}

	return createdReservation, nil
}

func (s *reservationService) UpdateStatus(ctx context.Context, ids []uint32, status string) error {
	atomic := func(txRepo postgresrepository.PostgresRepository) error {
		return txRepo.Reservation().UpdateStatus(ctx, ids, status)
	}

	err := s.Repo.Postgres().Atomic(ctx, s.Config, atomic)
	if err != nil {
		return err
	}

	return nil
}
