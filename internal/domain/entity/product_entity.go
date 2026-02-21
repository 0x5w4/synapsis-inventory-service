package entity

type Product struct {
	Base
	ID    uint32
	Name  string
	Stock int
	Price float64
}
