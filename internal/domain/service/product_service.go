package service

import (
	"context"
	"inventory-service/config"
	"inventory-service/internal/adapter/repository"
	postgresrepository "inventory-service/internal/adapter/repository/postgres"
	"inventory-service/internal/domain/entity"
	"inventory-service/pkg/logger"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	DeleteProduct(ctx context.Context, id uint32) (string, error)
	Find(ctx context.Context, filter *postgresrepository.FilterProductPayload) ([]*entity.Product, int, error)
	FindByID(ctx context.Context, id uint) (*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
}

type productService struct {
	config *config.Config
	repo   repository.Repository
	logger logger.Logger
}

func NewProductService(config *config.Config, repo repository.Repository, logger logger.Logger) *productService {
	return &productService{config: config, repo: repo, logger: logger}
}

func (s *productService) Find(ctx context.Context, filter *postgresrepository.FilterProductPayload) ([]*entity.Product, int, error) {
	return s.repo.Postgres().Product().Find(ctx, filter)
}

func (s *productService) FindByID(ctx context.Context, id uint) (*entity.Product, error) {
	return s.repo.Postgres().Product().FindByID(ctx, id)
}

func (s *productService) CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	// Implement logic for creating a product
	return product, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uint32) (string, error) {
	// Implement logic for deleting a product
	return "Product deleted successfully", nil
}

func (s *productService) UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	// Implement logic for updating a product
	return product, nil
}
