package entity

type Reservation struct {
	Base

	ProductID uint32
	OrderID   uint32
	Quantity  int
	Status    string

	Product *Product
}
