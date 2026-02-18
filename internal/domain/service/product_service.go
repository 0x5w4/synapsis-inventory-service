package service

import (
    "context"
    "goapptemp/config"
    "goapptemp/internal/adapter/repository/postgres"
    "goapptemp/internal/adapter/repository"
    "goapptemp/internal/domain/entity"
    "goapptemp/pkg/logger"
)

type ProductService interface {
    Find(ctx context.Context, filter *postgres.FilterProductPayload) ([]*entity.Product, int, error)
    FindByID(ctx context.Context, id uint) (*entity.Product, error)
}

type productService struct {
    config *config.Config
    repo   repository.Repository
    logger logger.Logger
}

func NewProductService(config *config.Config, repo repository.Repository, logger logger.Logger) *productService {
    return &productService{config: config, repo: repo, logger: logger}
}

func (s *productService) Find(ctx context.Context, filter *postgres.FilterProductPayload) ([]*entity.Product, int, error) {
    return s.repo.Postgres().Product().Find(ctx, filter)
}

func (s *productService) FindByID(ctx context.Context, id uint) (*entity.Product, error) {
    return s.repo.Postgres().Product().FindByID(ctx, id)
}
