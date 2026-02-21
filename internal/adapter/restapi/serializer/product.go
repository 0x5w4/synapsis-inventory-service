package serializer

import (
	"inventory-service/internal/domain/entity"
	"time"
)

type ProductResponse struct {
	ID        uint32    `json:"id"`
	Name      string    `json:"name"`
	Stock     int       `json:"stock"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func SerializeProduct(arg *entity.Product) *ProductResponse {
	if arg == nil {
		return nil
	}

	return &ProductResponse{
		ID:        arg.ID,
		Name:      arg.Name,
		Stock:     arg.Stock,
		Price:     arg.Price,
		CreatedAt: arg.CreatedAt,
		UpdatedAt: arg.UpdatedAt,
	}
}

func SerializeProducts(arg []*entity.Product) []*ProductResponse {
	if len(arg) == 0 {
		return nil
	}

	res := make([]*ProductResponse, 0, len(arg))

	for i := range arg {
		if arg[i] == nil {
			continue
		}

		res = append(res, SerializeProduct(arg[i]))
	}

	return res
}
