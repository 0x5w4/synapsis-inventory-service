package service

import (
    "context"
    "goapptemp/config"
    "goapptemp/internal/adapter/repository"
    "goapptemp/internal/adapter/repository/postgres"
    "goapptemp/internal/domain/entity"
    "goapptemp/pkg/logger"
)

type ReservationService interface {
    Create(ctx context.Context, reservation *entity.Reservation) (*entity.Reservation, error)
    Find(ctx context.Context, filter *postgres.FilterReservationPayload) ([]*entity.Reservation, int, error)
    FindByID(ctx context.Context, id uint) (*entity.Reservation, error)
}

type reservationService struct {
    config *config.Config
    repo   repository.Repository
    logger logger.Logger
}

func NewReservationService(config *config.Config, repo repository.Repository, logger logger.Logger) *reservationService {
    return &reservationService{config: config, repo: repo, logger: logger}
}

func (s *reservationService) Create(ctx context.Context, reservation *entity.Reservation) (*entity.Reservation, error) {
    // Use Postgres repository to create reservation in a transaction
    var res *entity.Reservation

    atomic := func(txRepo postgres.PostgresRepository) error {
        var err error
        res, err = txRepo.Reservation().Create(ctx, reservation)
        return err
    }

    if err := s.repo.Postgres().Atomic(ctx, s.config, atomic); err != nil {
        return nil, err
    }

    return res, nil
}

func (s *reservationService) Find(ctx context.Context, filter *postgres.FilterReservationPayload) ([]*entity.Reservation, int, error) {
    return s.repo.Postgres().Reservation().Find(ctx, filter)
}

func (s *reservationService) FindByID(ctx context.Context, id uint) (*entity.Reservation, error) {
    return s.repo.Postgres().Reservation().FindByID(ctx, id)
}
