package models

type Payment struct {
	BaseGormModel
	Method string  `json:"method" gorm:"not null"`
	Orders []Order `json:"orders" gorm:"foreignKey:PaymentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
