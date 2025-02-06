package models

type Payment struct {
	BaseGormModel
	Method string `json:"method" gorm:"not null"`
}
