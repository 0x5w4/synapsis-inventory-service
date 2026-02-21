package service

import (
	"inventory-service/config"
	"inventory-service/internal/adapter/repository"
	"inventory-service/pkg/logger"
)

type Properties struct {
	Config *config.Config
	Repo   repository.Repository
	Logger logger.Logger
}
