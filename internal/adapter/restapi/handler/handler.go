package handler

import (
	"errors"
	"inventory-service/config"
	"inventory-service/internal/domain/service"
	"inventory-service/pkg/logger"

	validator "github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
)

type Handler interface {
	Product() ProductHandler
}

type properties struct {
	config    *config.Config
	logger    logger.Logger
	service   service.Service
	validator *validator.Validate
	db        *bun.DB
}

type handler struct {
	properties
	productHandler ProductHandler
}

func NewHandler(config *config.Config, logger logger.Logger, service service.Service, db *bun.DB) (*handler, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	validate := validator.New()

	props := properties{
		config:    config,
		service:   service,
		logger:    logger,
		validator: validate,
		db:        db,
	}

	h := &handler{
		properties:     props,
		productHandler: NewProductHandler(props),
	}

	return h, nil
}

func (h *handler) Product() ProductHandler {
	return h.productHandler
}
