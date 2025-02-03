package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/tuantran0910/rainbow/internal/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	WithTX(tx *gorm.DB) IUserRepository
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) WithTX(tx *gorm.DB) IUserRepository {
	if tx == nil {
		return ur
	}

	return &userRepository{
		db: tx,
	}
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := ur.db.WithContext(ctx).Take(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

func (ur *userRepository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	var user models.User
	if err := ur.db.WithContext(ctx).Take(&user, "phone_number = ?", phoneNumber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

func (ur *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	if err := ur.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
