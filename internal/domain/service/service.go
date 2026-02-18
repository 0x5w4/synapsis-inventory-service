package service

import (
	"context"
	"inventory-service/config"
	"inventory-service/internal/adapter/repository"
	"inventory-service/pkg/logger"
	"inventory-service/proto/pb"
)

var _ Service = (*service)(nil)
var _ pb.InventoryServiceServer = (*service)(nil)

type Service interface {
	ProductService
	ReservationService
	mustEmbedUnimplementedInventoryServiceServer()
}

type service struct {
	productService     ProductService
	reservationService ReservationService
}

func NewService(
	config *config.Config,
	repo repository.Repository,
	logger logger.Logger,
) (*service, error) {
	return &service{
		productService:     NewProductService(config, repo, logger),
		reservationService: NewReservationService(config, repo, logger),
	}, nil
}

func (s *service) Product() ProductService {
	return s.productService
}

func (s *service) Reservation() ReservationService {
	return s.reservationService
}

func (s *service) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	return &pb.ListProductsResponse{}, nil
}

func (s *service) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	return &pb.GetProductResponse{}, nil
}

func (s *service) CreateReservation(ctx context.Context, req *pb.CreateReservationRequest) (*pb.CreateReservationResponse, error) {
	// Implement logic for creating a reservation
	return &pb.CreateReservationResponse{}, nil
}

func (s *service) ListReservations(ctx context.Context, req *pb.ListReservationsRequest) (*pb.ListReservationsResponse, error) {
	// Implement logic for listing reservations
	return &pb.ListReservationsResponse{}, nil
}

func (s *service) GetReservation(ctx context.Context, req *pb.GetReservationRequest) (*pb.GetReservationResponse, error) {
	// Implement logic for getting a reservation
	return &pb.GetReservationResponse{}, nil
}

func (s *service) mustEmbedUnimplementedInventoryServiceServer() {}
