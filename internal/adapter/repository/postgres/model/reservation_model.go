package model

import (
	"inventory-service/internal/domain/entity"

	"github.com/uptrace/bun"
)

type Reservation struct {
	bun.BaseModel `bun:"table:reservations,alias:reservation"`
	Base
	ProductID uint32 `bun:"product_id,notnull"`
	OrderID   uint32 `bun:"order_id,notnull"`
	Quantity  int    `bun:"quantity,notnull"`
	Status    string `bun:"status,notnull"`

	Product *Product `bun:"rel:belongs-to,join:product_id=id"`
}

func (m *Reservation) ToDomain() *entity.Reservation {
	if m == nil {
		return nil
	}

	return &entity.Reservation{
		Base: entity.Base{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: m.DeletedAt,
		},
		ProductID: m.ProductID,
		OrderID:   m.OrderID,
		Quantity:  m.Quantity,
		Status:    m.Status,
		Product:   m.Product.ToDomain(),
	}
}

func ToReservationsDomain(arg []*Reservation) []*entity.Reservation {
	if len(arg) == 0 {
		return nil
	}

	res := make([]*entity.Reservation, 0, len(arg))

	for i := range arg {
		if arg[i] == nil {
			continue
		}

		res = append(res, arg[i].ToDomain())
	}

	return res
}

func AsReservation(arg *entity.Reservation) *Reservation {
	if arg == nil {
		return nil
	}

	return &Reservation{
		Base: Base{
			ID:        arg.ID,
			CreatedAt: arg.CreatedAt,
			UpdatedAt: arg.UpdatedAt,
			DeletedAt: arg.DeletedAt,
		},
		ProductID: arg.ProductID,
		OrderID:   arg.OrderID,
		Quantity:  arg.Quantity,
		Status:    arg.Status,
		Product:   AsProduct(arg.Product),
	}
}

func AsReservations(arg []*entity.Reservation) []*Reservation {
	if len(arg) == 0 {
		return nil
	}

	res := make([]*Reservation, 0, len(arg))

	for i := range arg {
		if arg[i] == nil {
			continue
		}

		res = append(res, AsReservation(arg[i]))
	}

	return res
}
