package service

import (
	"context"
	"inventory-service/config"
	"inventory-service/internal/adapter/repository"
	postgresrepository "inventory-service/internal/adapter/repository/postgres"
	"inventory-service/internal/domain/entity"
	"inventory-service/pkg/logger"
	"inventory-service/proto/pb"
)

type ReservationService interface {
	Create(ctx context.Context, reservation *entity.Reservation) (*entity.Reservation, error)
	Find(ctx context.Context, filter *pb.FilterReservationPayload) ([]*pb.Reservation, int, error)
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

	atomic := func(txRepo postgresrepository.PostgresRepository) error {
		var err error
		res, err = txRepo.Reservation().Create(ctx, reservation)
		return err
	}

	if err := s.repo.Postgres().Atomic(ctx, s.config, atomic); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *reservationService) Find(ctx context.Context, filter *pb.FilterReservationPayload) ([]*pb.Reservation, int, error) {
	// Use pb.FilterReservationPayload directly
	reservations, total, err := s.repo.Reservation().Find(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Map domain reservations to pb.Reservation
	pbReservations := make([]*pb.Reservation, len(reservations))
	for i, reservation := range reservations {
		pbReservations[i] = &pb.Reservation{
			Id:        reservation.ID,
			ProductId: reservation.ProductID,
			OrderId:   reservation.OrderID,
			Quantity:  reservation.Quantity,
			Status:    reservation.Status,
		}
	}

	return pbReservations, total, nil
}

func (s *reservationService) FindByID(ctx context.Context, id uint) (*entity.Reservation, error) {
	return s.repo.Postgres().Reservation().FindByID(ctx, id)
}
