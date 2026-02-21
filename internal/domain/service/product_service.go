package service

import (
	"context"
	postgresrepository "inventory-service/internal/adapter/repository/postgres"
	"inventory-service/internal/domain/entity"
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
	Properties
}

func NewProductService(props Properties) *productService {
	return &productService{Properties: props}
}

func (s *productService) Find(ctx context.Context, filter *postgresrepository.FilterProductPayload) ([]*entity.Product, int, error) {
	return s.Repo.Postgres().Product().Find(ctx, filter)
}

func (s *productService) FindByID(ctx context.Context, id uint32) (*entity.Product, error) {
	return s.Repo.Postgres().Product().FindByID(ctx, id)
}

func (s *productService) Create(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	var createdProduct *entity.Product

	atomic := func(r postgresrepository.PostgresRepository) error {
		var err error
		createdProduct, err = r.Product().Create(ctx, product)
		return err
	}

	err := s.Repo.Postgres().Atomic(ctx, s.Config, atomic)
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

	err := s.Repo.Postgres().Atomic(ctx, s.Config, atomic)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (s *productService) Delete(ctx context.Context, id uint32) error {
	atomic := func(r postgresrepository.PostgresRepository) error {
		return r.Product().Delete(ctx, id)
	}

	err := s.Repo.Postgres().Atomic(ctx, s.Config, atomic)
	if err != nil {
		return err
	}

	return nil
}
