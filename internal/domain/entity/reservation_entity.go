package entity

type Reservation struct {
	Base
	ProductID        uint
	OrderID        uint
	Quantity       int
	Status         string
}