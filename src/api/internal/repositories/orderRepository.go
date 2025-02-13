package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/models"
	"github.com/tuantran0910/rainbow/pkg/pagination"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	WithTX(tx *gorm.DB) IOrderRepository
	GetOrdersByUserId(
		ctx context.Context,
		page, limit int,
		userId uuid.UUID,
	) ([]*models.Order, *pagination.Pagination, error)
	GetOrderById(ctx context.Context, orderId uuid.UUID, userId uuid.UUID) (*models.Order, error)
	CreateOrder(ctx context.Context, order *models.Order) error
	CreateOrderItem(ctx context.Context, orderItem *models.OrderItem) error
	UpdateOrder(ctx context.Context, order *models.Order) error
	DeleteOrder(ctx context.Context, orderId uuid.UUID) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (or *orderRepository) WithTX(tx *gorm.DB) IOrderRepository {
	if tx == nil {
		return or
	}
	return &orderRepository{
		db: tx,
	}
}

func (or *orderRepository) GetOrdersByUserId(
	ctx context.Context,
	page, limit int,
	userId uuid.UUID,
) ([]*models.Order, *pagination.Pagination, error) {
	var totalOrders int64
	if err := or.db.WithContext(ctx).Model(&models.Order{}).Count(&totalOrders).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch total number of orders: %w", err)
	}

	pagination := pagination.NewPagination(page, limit, int(totalOrders))

	var orders []*models.Order
	err := or.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Offset(pagination.Offset).
		Limit(limit).
		Find(&orders).Error

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	return orders, pagination, nil
}

func (or *orderRepository) GetOrderById(
	ctx context.Context,
	orderId uuid.UUID,
	userId uuid.UUID,
) (*models.Order, error) {
	var order *models.Order
	if err := or.db.WithContext(ctx).Preload("OrderItems").Take(&order, "id = ? AND user_id = ?", orderId, userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}
	return order, nil
}

func (or *orderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	if err := or.db.WithContext(ctx).Create(order).Error; err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

func (or *orderRepository) CreateOrderItem(ctx context.Context, orderItem *models.OrderItem) error {
	if err := or.db.WithContext(ctx).Create(orderItem).Error; err != nil {
		return fmt.Errorf("failed to create order item: %w", err)
	}
	return nil
}

func (or *orderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	if err := or.db.WithContext(ctx).Save(order).Error; err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	return nil
}

func (or *orderRepository) DeleteOrder(ctx context.Context, orderId uuid.UUID) error {
	if err := or.db.WithContext(ctx).Unscoped().Where("id = ?", orderId).Delete(&models.Order{}).Error; err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}
