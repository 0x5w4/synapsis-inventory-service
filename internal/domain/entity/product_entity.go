package entity

type Product struct {
	Base
	Code        string
	Name        string
	Description *string
	Stock       int
}