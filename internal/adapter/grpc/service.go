package grpcadapter

import (
	"context"
	"inventory-service/config"
	"inventory-service/internal/adapter/repository"
	postgresrepository "inventory-service/internal/adapter/repository/postgres"
	"inventory-service/internal/domain/entity"
	"inventory-service/internal/domain/service"
	"inventory-service/pkg/logger"
	"inventory-service/proto/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

var _ pb.InventoryServiceServer = (*grpcService)(nil)

type grpcService struct {
	pb.UnimplementedInventoryServiceServer
	productService     service.ProductService
	reservationService service.ReservationService
}

func NewGRPCService(
	config *config.Config,
	repo repository.Repository,
	logger logger.Logger,
) (*grpcService, error) {
	return &grpcService{
		productService:     service.NewProductService(config, repo, logger),
		reservationService: service.NewReservationService(config, repo, logger),
	}, nil
}

func (s *grpcService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	filter := &postgresrepository.FilterProductPayload{
		IDs:     req.Ids,
		Names:   req.Names,
		Search:  req.Search,
		Page:    int(req.Page),
		PerPage: int(req.PerPage),
	}

	products, total, err := s.productService.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := &pb.ListProductsResponse{
		Total:    int32(total),
		Products: make([]*pb.Product, len(products)),
	}

	for i, product := range products {
		response.Products[i] = &pb.Product{
			Id:    product.Base.ID,
			Name:  product.Name,
			Stock: int32(product.Stock),
		}
	}

	return response, nil
}

func (s *grpcService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	product, err := s.productService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Product{
		Id:    product.Base.ID,
		Name:  product.Name,
		Stock: int32(product.Stock),
	}, nil
}

func (s *grpcService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	productEntity := &entity.Product{
		Name:  req.Name,
		Stock: int(req.Stock),
	}

	createdProduct, err := s.productService.Create(ctx, productEntity)
	if err != nil {
		return nil, err
	}

	response := &pb.Product{
		Id:    createdProduct.Base.ID,
		Name:  createdProduct.Name,
		Stock: int32(createdProduct.Stock),
	}

	return response, nil
}

func (s *grpcService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	product := &entity.Product{
		Base:  entity.Base{ID: req.Id},
		Name:  req.Name,
		Stock: int(req.Stock),
	}

	updatedProduct, err := s.productService.Update(ctx, product)
	if err != nil {
		return nil, err
	}

	return &pb.Product{
		Id:    updatedProduct.Base.ID,
		Name:  updatedProduct.Name,
		Stock: int32(updatedProduct.Stock),
	}, nil
}

func (s *grpcService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	if err := s.productService.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *grpcService) CreateReservation(ctx context.Context, req *pb.CreateReservationRequest) (*pb.Reservation, error) {
	reservation := &entity.Reservation{
		ProductID: req.ProductId,
		OrderID:   req.OrderId,
		Quantity:  int(req.Quantity),
	}

	createdReservation, err := s.reservationService.Create(ctx, reservation)
	if err != nil {
		return nil, err
	}

	return &pb.Reservation{
		Id:        createdReservation.Base.ID,
		ProductId: createdReservation.ProductID,
		OrderId:   createdReservation.OrderID,
		Quantity:  int32(createdReservation.Quantity),
		Status:    MapDBStatusToPBStatus(createdReservation.Status),
	}, nil
}

func (s *grpcService) ListReservations(ctx context.Context, req *pb.ListReservationsRequest) (*pb.ListReservationsResponse, error) {
	filter := &postgresrepository.FilterReservationPayload{
		ProductIDs: req.ProductIds,
		OrderIDs:   req.OrderIds,
		Page:       int(req.Page),
		PerPage:    int(req.PerPage),
	}

	filter.Statuses = make([]string, len(req.Statuses))
	for i, status := range req.Statuses {
		filter.Statuses[i] = MapPBStatusToDBStatus(status)
	}

	reservations, total, err := s.reservationService.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := &pb.ListReservationsResponse{
		Total:        int32(total),
		Reservations: make([]*pb.Reservation, len(reservations)),
	}

	for i, reservation := range reservations {
		response.Reservations[i] = &pb.Reservation{
			Id:        reservation.Base.ID,
			ProductId: reservation.ProductID,
			OrderId:   reservation.OrderID,
			Quantity:  int32(reservation.Quantity),
			Status:    MapDBStatusToPBStatus(reservation.Status),
		}
	}

	return response, nil
}

func (s *grpcService) GetReservation(ctx context.Context, req *pb.GetReservationRequest) (*pb.Reservation, error) {
	reservation, err := s.reservationService.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Reservation{
		Id:        reservation.Base.ID,
		ProductId: reservation.ProductID,
		OrderId:   reservation.OrderID,
		Quantity:  int32(reservation.Quantity),
		Status:    MapDBStatusToPBStatus(reservation.Status),
	}, nil
}

func (s *grpcService) mustEmbedUnimplementedInventoryServiceServer() {}
