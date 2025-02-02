package models

type Product struct {
	BaseGormModel
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
