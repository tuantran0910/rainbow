package models

import "github.com/google/uuid"

type Order struct {
	BaseGormModel
	UserID          uuid.UUID    `json:"user_id"          gorm:"not null"`
	PaymentID       uuid.UUID    `json:"payment_id"       gorm:"not null"`
	PromotionID     *uuid.UUID   `json:"promotion_id"                                                                            gore:"null"`
	ShippingAddress string       `json:"shipping_address" gorm:"not null"`
	TotalAmount     float64      `json:"total_amount"     gorm:"not null;check:total_amount >= 0"`
	OrderItems      []*OrderItem `json:"order_items"      gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type OrderItem struct {
	BaseGormModel
	OrderID   uuid.UUID `json:"order_id"   gorm:"not null"`
	BookID    uuid.UUID `json:"book_id"    gorm:"not null"`
	Quantity  int       `json:"quantity"   gorm:"not null;check:quantity > 0"`
	UnitPrice float64   `json:"unit_price" gorm:"not null;check:unit_price >= 0"`
	Discount  float64   `json:"discount"   gorm:"not null;check:discount >= 0"`
}
