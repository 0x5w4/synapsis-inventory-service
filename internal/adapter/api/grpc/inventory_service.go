package grpcapi

import (
    "context"
    "time"
    "goapptemp/internal/adapter/api/grpc/pb"
    "goapptemp/internal/domain/service"
    "goapptemp/pkg/logger"
    "google.golang.org/grpc"
)

type InventoryServer struct {
    pb.UnimplementedInventoryServiceServer
    svc service.Service
    logger logger.Logger
}

func NewInventoryServer(svc service.Service, logger logger.Logger) *InventoryServer {
    return &InventoryServer{svc: svc, logger: logger}
}

func (s *InventoryServer) Register(grpcServer *grpc.Server) {
    pb.RegisterInventoryServiceServer(grpcServer, s)
}

func (s *InventoryServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
    filter := &service.PostgresFilterProductPayload{}
    // Map pagination/search
    prodFilter := &service.PostgresFilterProductPayload{}
    products, total, err := s.svc.Product().Find(ctx, prodFilter)
    if err != nil {
        return nil, err
    }

    res := &pb.ListProductsResponse{Total: int32(total)}
    for _, p := range products {
        res.Products = append(res.Products, &pb.Product{
            Id: uint32(p.ID),
            Code: p.Code,
            Name: p.Name,
            Description: derefString(p.Description),
            Stock: int32(p.Stock),
            CreatedAt: p.CreatedAt,
            UpdatedAt: p.UpdatedAt,
        })
    }

    return res, nil
}

func (s *InventoryServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
    prod, err := s.svc.Product().FindByID(ctx, uint(req.Id))
    if err != nil {
        return nil, err
    }

    res := &pb.GetProductResponse{Product: &pb.Product{
        Id: uint32(prod.ID), Code: prod.Code, Name: prod.Name,
        Description: derefString(prod.Description), Stock: int32(prod.Stock), CreatedAt: prod.CreatedAt, UpdatedAt: prod.UpdatedAt,
    }}

    return res, nil
}

func (s *InventoryServer) CreateReservation(ctx context.Context, req *pb.CreateReservationRequest) (*pb.CreateReservationResponse, error) {
    r := req.Reservation
    entityRes := &service.ReservationEntity{ProductID: uint(r.ProductId), OrderID: uint(r.OrderId), Quantity: int(r.Quantity), Status: r.Status}
    created, err := s.svc.Reservation().Create(ctx, entityRes)
    if err != nil {
        return nil, err
    }

    return &pb.CreateReservationResponse{Reservation: &pb.Reservation{Id: uint32(created.ID), ProductId: uint32(created.ProductID), OrderId: uint32(created.OrderID), Quantity: int32(created.Quantity), Status: created.Status, CreatedAt: created.CreatedAt, UpdatedAt: created.UpdatedAt}}, nil
}

func (s *InventoryServer) ListReservations(ctx context.Context, req *pb.ListReservationsRequest) (*pb.ListReservationsResponse, error) {
    filter := &service.PostgresFilterReservationPayload{}
    reservations, total, err := s.svc.Reservation().Find(ctx, filter)
    if err != nil {
        return nil, err
    }

    res := &pb.ListReservationsResponse{Total: int32(total)}
    for _, r := range reservations {
        res.Reservations = append(res.Reservations, &pb.Reservation{Id: uint32(r.ID), ProductId: uint32(r.ProductID), OrderId: uint32(r.OrderID), Quantity: int32(r.Quantity), Status: r.Status, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt})
    }

    return res, nil
}

func (s *InventoryServer) GetReservation(ctx context.Context, req *pb.GetReservationRequest) (*pb.GetReservationResponse, error) {
    r, err := s.svc.Reservation().FindByID(ctx, uint(req.Id))
    if err != nil {
        return nil, err
    }

    return &pb.GetReservationResponse{Reservation: &pb.Reservation{Id: uint32(r.ID), ProductId: uint32(r.ProductID), OrderId: uint32(r.OrderID), Quantity: int32(r.Quantity), Status: r.Status, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt}}, nil
}

func derefString(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}
