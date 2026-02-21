package handler

import (
	postgresrepository "inventory-service/internal/adapter/repository/postgres"
	"inventory-service/internal/adapter/restapi/response"
	"inventory-service/internal/adapter/restapi/serializer"
	"inventory-service/internal/domain/entity"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProductHandler interface {
	Create(c echo.Context) error
	Get(c echo.Context) error
	List(c echo.Context) error
	Update(c echo.Context) error
}

type productHandler struct {
	properties
}

func NewProductHandler(props properties) ProductHandler {
	return &productHandler{properties: props}
}

type CreateProductRequest struct {
	Name  string  `json:"name" validate:"required"`
	Stock int     `json:"stock" validate:"required,min=0"`
	Price float64 `json:"price" validate:"required,min=0"`
}

func (h *productHandler) Create(c echo.Context) error {
	var req CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}

	product := &entity.Product{
		Name:  req.Name,
		Stock: req.Stock,
		Price: req.Price,
	}

	createdProduct, err := h.service.Product().Create(c.Request().Context(), product)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, serializer.SerializeProduct(createdProduct))
}

func (h *productHandler) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return err
	}

	product, err := h.service.Product().FindByID(c.Request().Context(), uint32(id))
	if err != nil {
		return err
	}

	return response.Success(c, "Product retrieved successfully", serializer.SerializeProduct(product))
}

func (h *productHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))

	filter := &postgresrepository.FilterProductPayload{
		Page:    page,
		PerPage: perPage,
	}

	products, total, err := h.service.Product().Find(c.Request().Context(), filter)
	if err != nil {
		return err
	}

	return response.Paginate(c, "Products retrieved successfully", serializer.SerializeProducts(products), response.Pagination{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		TotalPage:  (total + perPage - 1) / perPage,
	})
}

func (h *productHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return err
	}

	var req CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := h.validator.Struct(req); err != nil {
		return err
	}

	product := &entity.Product{
		ID:    uint32(id),
		Name:  req.Name,
		Stock: req.Stock,
		Price: req.Price,
	}

	updatedProduct, err := h.service.Product().Update(c.Request().Context(), product)
	if err != nil {
		return err
	}

	return response.Success(c, "Product updated successfully", serializer.SerializeProduct(updatedProduct))
}
