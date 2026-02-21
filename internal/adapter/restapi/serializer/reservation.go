package serializer

import (
	"inventory-service/internal/domain/entity"
	"time"
)

type ReservationResponse struct {
	ID        uint32           `json:"id"`
	ProductID uint32           `json:"product_id"`
	OrderID   uint32           `json:"order_id"`
	Quantity  int              `json:"quantity"`
	Status    string           `json:"status"`
	Product   *ProductResponse `json:"product"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func SerializeReservation(arg *entity.Reservation) *ReservationResponse {
	if arg == nil {
		return nil
	}

	return &ReservationResponse{
		ID:        arg.ID,
		ProductID: arg.ProductID,
		OrderID:   arg.OrderID,
		Quantity:  arg.Quantity,
		Status:    arg.Status,
		Product:   SerializeProduct(arg.Product),
		CreatedAt: arg.CreatedAt,
		UpdatedAt: arg.UpdatedAt,
	}
}

func SerializeReservations(arg []*entity.Reservation) []*ReservationResponse {
	if len(arg) == 0 {
		return nil
	}

	res := make([]*ReservationResponse, 0, len(arg))

	for i := range arg {
		if arg[i] == nil {
			continue
		}

		res = append(res, SerializeReservation(arg[i]))
	}

	return res
}
