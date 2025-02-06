package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/models"
	"gorm.io/gorm"
)

type IInventoryRepository interface {
	WithTX(tx *gorm.DB) IInventoryRepository
	GetInventoryByBookID(ctx context.Context, bookID uuid.UUID) (*models.Inventory, error)
	CreateInventory(ctx context.Context, inventory *models.Inventory) error
	UpdateInventory(ctx context.Context, inventoryId uuid.UUID, inventory *models.Inventory) error
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) IInventoryRepository {
	return &inventoryRepository{
		db: db,
	}
}

func (ir *inventoryRepository) WithTX(tx *gorm.DB) IInventoryRepository {
	if tx == nil {
		return ir
	}
	return &inventoryRepository{
		db: tx,
	}
}

func (ir *inventoryRepository) GetInventoryByBookID(ctx context.Context, bookID uuid.UUID) (*models.Inventory, error) {
	var inventory models.Inventory
	if err := ir.db.WithContext(ctx).Take(&inventory, "book_id = ?", bookID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch inventory with id %s: %w", bookID, err)
	}
	return &inventory, nil
}

func (ir *inventoryRepository) CreateInventory(ctx context.Context, inventory *models.Inventory) error {
	if err := ir.db.WithContext(ctx).Create(inventory).Error; err != nil {
		return fmt.Errorf("failed to create inventory: %w", err)
	}
	return nil
}

func (ir *inventoryRepository) UpdateInventory(ctx context.Context, inventoryId uuid.UUID, inventory *models.Inventory) error {
	if err := ir.db.WithContext(ctx).Model(&models.Inventory{}).Where("id = ?", inventoryId).Updates(inventory).Error; err != nil {
		return fmt.Errorf("failed to update inventory with id %s: %w", inventoryId, err)
	}
	return nil
}
