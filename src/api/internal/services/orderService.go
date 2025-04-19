package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/dtos"
	"github.com/tuantran0910/rainbow/internal/models"
	"github.com/tuantran0910/rainbow/internal/repositories"
	"github.com/tuantran0910/rainbow/pkg/pagination"
	"gorm.io/gorm"
)

type IOrderService interface {
	GetOrdersByUserId(
		ctx context.Context,
		page, limit int,
		currentUserId uuid.UUID,
	) ([]*models.Order, *pagination.Pagination, error)
	GetOrderById(
		ctx context.Context,
		orderId uuid.UUID,
		currentUserId uuid.UUID,
	) (*models.Order, error)
	CreateOrder(
		ctx context.Context,
		orderRequest dtos.CreateOrderRequest,
		currentUserId uuid.UUID,
	) (*models.Order, error)
	DeleteOrder(ctx context.Context, orderId uuid.UUID, currentUserId uuid.UUID) error
}

type orderService struct {
	db                      *gorm.DB
	orderRepository         repositories.IOrderRepository
	bookRepository          repositories.IBookRepository
	inventoryRepository     repositories.IInventoryRepository
	promotionRepository     repositories.IPromotionRepository
	userRepository          repositories.IUserRepository
	userPromotionRepository repositories.IUserPromotionRepository
}

func NewOrderService(
	db *gorm.DB,
	orderRepository repositories.IOrderRepository,
	bookRepository repositories.IBookRepository,
	inventoryRepository repositories.IInventoryRepository,
	promoRepository repositories.IPromotionRepository,
	userRepository repositories.IUserRepository,
	userPromotionRepository repositories.IUserPromotionRepository,
) IOrderService {
	return &orderService{
		db:                      db,
		orderRepository:         orderRepository,
		bookRepository:          bookRepository,
		inventoryRepository:     inventoryRepository,
		promotionRepository:     promoRepository,
		userRepository:          userRepository,
		userPromotionRepository: userPromotionRepository,
	}
}

func (os *orderService) withTX(
	ctx context.Context,
	fn func(
		context.Context,
		repositories.IOrderRepository,
		repositories.IBookRepository,
		repositories.IInventoryRepository,
		repositories.IPromotionRepository,
		repositories.IUserRepository,
		repositories.IUserPromotionRepository,
	) error,
) error {
	return os.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		orderRepository := os.orderRepository.WithTX(tx)
		bookRepository := os.bookRepository.WithTX(tx)
		inventoryRepository := os.inventoryRepository.WithTX(tx)
		promotionRepository := os.promotionRepository.WithTX(tx)
		userRepository := os.userRepository.WithTX(tx)
		userPromotionRepository := os.userPromotionRepository.WithTX(tx)
		return fn(
			ctx,
			orderRepository,
			bookRepository,
			inventoryRepository,
			promotionRepository,
			userRepository,
			userPromotionRepository,
		)
	})
}

func (os *orderService) GetOrdersByUserId(
	ctx context.Context,
	page, limit int,
	currentUserId uuid.UUID,
) ([]*models.Order, *pagination.Pagination, error) {
	return os.orderRepository.GetOrdersByUserId(ctx, page, limit, currentUserId)
}

func (os *orderService) GetOrderById(
	ctx context.Context,
	orderId uuid.UUID,
	currentUserId uuid.UUID,
) (*models.Order, error) {
	return os.orderRepository.GetOrderById(ctx, orderId, currentUserId)
}

func (os *orderService) CreateOrder(
	ctx context.Context,
	orderRequest dtos.CreateOrderRequest,
	currentUserId uuid.UUID,
) (*models.Order, error) {
	var createdOrder *models.Order
	err := os.withTX(ctx, func(
		ctx context.Context,
		orderRepository repositories.IOrderRepository,
		bookRepository repositories.IBookRepository,
		inventoryRepository repositories.IInventoryRepository,
		promotionRepository repositories.IPromotionRepository,
		userRepository repositories.IUserRepository,
		userPromotionRepository repositories.IUserPromotionRepository,
	) error {
		var err error

		order := &models.Order{
			UserID:          currentUserId,
			PaymentID:       orderRequest.PaymentId,
			PromotionID:     orderRequest.PromotionId,
			ShippingAddress: orderRequest.ShippingAddress,
		}
		if err = orderRepository.CreateOrder(ctx, order); err != nil {
			return err
		}

		bookTotalAmount := 0.0
		for _, orderItem := range orderRequest.OrderItems {
			book, err := bookRepository.GetBookById(ctx, orderItem.BookId, false)
			if err != nil {
				return err
			}
			if book == nil {
				return fmt.Errorf("book not found")
			}

			if book.Inventory.Stock < orderItem.Quantity {
				return fmt.Errorf("not enough stock for book with ID %s", book.ID)
			}

			bookTotalAmount += book.Price * float64(orderItem.Quantity)
			orderItemModel := &models.OrderItem{
				OrderID:   order.ID,
				BookID:    orderItem.BookId,
				Quantity:  orderItem.Quantity,
				UnitPrice: book.Price,
				Discount:  book.OriginalPrice - book.Price,
			}
			if err = orderRepository.CreateOrderItem(ctx, orderItemModel); err != nil {
				return err
			}

			bookInventory := book.Inventory
			bookInventory.Stock -= orderItem.Quantity
			if err = inventoryRepository.UpdateInventory(ctx, book.Inventory.ID, &bookInventory); err != nil {
				return err
			}

			book.SoldCount += orderItem.Quantity
			if err = bookRepository.UpdateBook(ctx, book.ID, book); err != nil {
				return err
			}
		}

		// Set the default total amount (without promotion)
		order.TotalAmount = bookTotalAmount

		var promotion *models.Promotion
		if orderRequest.PromotionId != nil {
			promotion, err = promotionRepository.GetPromotionById(ctx, *orderRequest.PromotionId)
			if err != nil {
				return err
			}

			// Check if the user has already used this promotion
			userPromotion, err := userPromotionRepository.GetUserPromotionByUserAndPromotion(
				ctx, currentUserId, *orderRequest.PromotionId,
			)
			if err != nil {
				return err
			}

			if userPromotion != nil {
				return fmt.Errorf("you have already used this promotion")
			}
		}

		if promotion != nil && promotion.UsedCount < promotion.MaxUses {
			if promotion.StartDate.After(order.CreatedAt) ||
				promotion.EndDate.Before(order.CreatedAt) {
				return fmt.Errorf("promotion is not available")
			}

			order.PromotionID = &promotion.ID
			if promotion.DiscountType == models.DiscountPercentage {
				order.TotalAmount = bookTotalAmount * (1 - promotion.DiscountValue/100)
			} else {
				order.TotalAmount = bookTotalAmount - promotion.DiscountValue
			}

			if order.TotalAmount < 0 {
				order.TotalAmount = 0
			}

			// Create a record in user_promotions table
			userPromotion := &models.UserPromotion{
				UserID:      currentUserId,
				PromotionID: promotion.ID,
			}
			if err = userPromotionRepository.CreateUserPromotion(ctx, userPromotion); err != nil {
				return err
			}

			promotion.UsedCount++
			if err = promotionRepository.UpdatePromotion(ctx, promotion.ID, promotion); err != nil {
				return err
			}
		}

		// Update the order with the final total amount
		if err = orderRepository.UpdateOrder(ctx, order); err != nil {
			return err
		}

		// Get the created order with all fields populated
		createdOrder, err = orderRepository.GetOrderById(ctx, order.ID, currentUserId)
		return err
	})
	return createdOrder, err
}

func (os *orderService) DeleteOrder(
	ctx context.Context,
	orderId uuid.UUID,
	currentUserId uuid.UUID,
) error {
	return os.withTX(ctx, func(
		ctx context.Context,
		orderRepository repositories.IOrderRepository,
		bookRepository repositories.IBookRepository,
		inventoryRepository repositories.IInventoryRepository,
		promotionRepository repositories.IPromotionRepository,
		userRepository repositories.IUserRepository,
		userPromotionRepository repositories.IUserPromotionRepository,
	) error {
		order, err := orderRepository.GetOrderById(ctx, orderId, currentUserId)
		if err != nil {
			return err
		}
		if order == nil {
			return fmt.Errorf("order not found")
		}

		for _, orderItem := range order.OrderItems {
			book, err := bookRepository.GetBookById(ctx, orderItem.BookID, false)
			if err != nil {
				return err
			}
			if book == nil {
				return fmt.Errorf("book not found")
			}

			book.Inventory.Stock += orderItem.Quantity
			if err := inventoryRepository.UpdateInventory(ctx, book.Inventory.ID, &book.Inventory); err != nil {
				return err
			}

			book.SoldCount -= orderItem.Quantity
			if err := bookRepository.UpdateBook(ctx, book.ID, book); err != nil {
				return err
			}
		}

		if order.PromotionID != nil {
			promotion, err := promotionRepository.GetPromotionById(ctx, *order.PromotionID)
			if err != nil {
				return err
			}
			if promotion != nil && promotion.StartDate.Before(order.CreatedAt) &&
				promotion.EndDate.After(order.CreatedAt) {
				promotion.UsedCount--
				if err := promotionRepository.UpdatePromotion(ctx, promotion.ID, promotion); err != nil {
					return err
				}
			}
		}

		return orderRepository.DeleteOrder(ctx, orderId)
	})
}
