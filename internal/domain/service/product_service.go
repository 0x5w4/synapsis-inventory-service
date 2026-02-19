package service

import (
	"context"
	"inventory-service/config"
	"inventory-service/internal/adapter/repository"
	postgresrepository "inventory-service/internal/adapter/repository/postgres"
	"inventory-service/internal/domain/entity"
	"inventory-service/pkg/logger"
)

var _ ProductService = (*productService)(nil)

type ProductService interface {
	Create(ctx context.Context, product *entity.Product) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) (*entity.Product, error)
	Delete(ctx context.Context, id uint32) error
	Find(ctx context.Context, filter *postgresrepository.FilterProductPayload) ([]*entity.Product, int, error)
	FindByID(ctx context.Context, id uint32) (*entity.Product, error)
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

func (s *productService) FindByID(ctx context.Context, id uint32) (*entity.Product, error) {
	return s.repo.Postgres().Product().FindByID(ctx, id)
}

func (s *productService) Create(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	var createdProduct *entity.Product

	atomic := func(r postgresrepository.PostgresRepository) error {
		var err error
		createdProduct, err = r.Product().Create(ctx, product)
		return err
	}

	err := s.repo.Postgres().Atomic(ctx, s.config, atomic)
	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}

func (s *productService) Update(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	var updatedProduct *entity.Product

	atomic := func(r postgresrepository.PostgresRepository) error {
		var err error
		updatedProduct, err = r.Product().Update(ctx, product)
		return err
	}

	err := s.repo.Postgres().Atomic(ctx, s.config, atomic)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (s *productService) Delete(ctx context.Context, id uint32) error {
	atomic := func(r postgresrepository.PostgresRepository) error {
		return r.Product().Delete(ctx, id)
	}

	err := s.repo.Postgres().Atomic(ctx, s.config, atomic)
	if err != nil {
		return err
	}

	return nil
}
