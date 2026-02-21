package service

import (
	"inventory-service/config"
	"inventory-service/internal/adapter/repository"
	"inventory-service/pkg/logger"
	"inventory-service/proto/pb"
)

var _ Service = (*service)(nil)

type Service interface {
	Product() ProductService
	Reservation() ReservationService
}

type Properties struct {
	Config                 *config.Config
	Repo                   repository.Repository
	Logger                 logger.Logger
	InventoryServiceClient pb.InventoryServiceClient
}

type service struct {
	Properties
	productService     ProductService
	reservationService ReservationService
}

func NewService(
	config *config.Config,
	repo repository.Repository,
	logger logger.Logger,
	inventoryServiceClient pb.InventoryServiceClient,
) (*service, error) {
	props := Properties{
		Config:                 config,
		Repo:                   repo,
		Logger:                 logger,
		InventoryServiceClient: inventoryServiceClient,
	}

	return &service{
		Properties:         props,
		productService:     NewProductService(props),
		reservationService: NewReservationService(props),
	}, nil
}

func (s *service) Product() ProductService {
	return s.productService
}

func (s *service) Reservation() ReservationService {
	return s.reservationService
}
