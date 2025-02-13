package dtos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateOrderRequest struct {
	PaymentId       uuid.UUID                `json:"payment_id"       binding:"required"`
	PromotionId     *uuid.UUID               `json:"promotion_id"`
	ShippingAddress string                   `json:"shipping_address" binding:"required"`
	OrderItems      []CreateOrderItemRequest `json:"order_items"      binding:"required"`
}

type CreateOrderItemRequest struct {
	BookId   uuid.UUID `json:"book_id"  binding:"required"`
	Quantity int       `json:"quantity" binding:"required,min=1"`
}

type GetOrderResponse struct {
	ID              uuid.UUID                `json:"id"`
	UserID          uuid.UUID                `json:"user_id"`
	PaymentID       uuid.UUID                `json:"payment_id"`
	ShippingAddress string                   `json:"shipping_address"`
	TotalAmount     float64                  `json:"total_amount"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	DeletedAt       gorm.DeletedAt           `json:"deleted_at,omitempty"`
	OrderItems      *[]*GetOrderItemResponse `json:"order_items"`
}

type GetOrderItemResponse struct {
	ID          uuid.UUID      `json:"id"`
	OrderID     uuid.UUID      `json:"order_id"`
	BookID      uuid.UUID      `json:"book_id"`
	PromotionID *uuid.UUID     `json:"promotion_id"`
	Quantity    int            `json:"quantity"`
	UnitPrice   float64        `json:"unit_price"`
	Discount    float64        `json:"discount"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type ListOrdersResponse struct {
	Orders []*GetOrderResponse `json:"orders"`
}
