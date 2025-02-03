package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/internal/models"
	"github.com/tuantran0910/rainbow/pkg/pagination"
	"gorm.io/gorm"
)

type IUserRepository interface {
	WithTX(tx *gorm.DB) IUserRepository
	GetUsers(ctx context.Context, page, limit int) ([]*models.User, *pagination.Pagination, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, userId uuid.UUID, user *models.User) error
	DeleteUser(ctx context.Context, userId uuid.UUID) error
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

func (ur *userRepository) GetUsers(ctx context.Context, page, limit int) ([]*models.User, *pagination.Pagination, error) {
	// Get total number of users
	var totalUsers int64
	if err := ur.db.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch total number of users: %w", err)
	}

	// Define pagination
	pagination := pagination.NewPagination(page, limit, int(totalUsers))

	var users []*models.User
	if err := ur.db.WithContext(ctx).Offset(pagination.Offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch users %w: ", err)
	}

	return users, pagination, nil
}

func (ur *userRepository) GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	var user models.User
	if err := ur.db.WithContext(ctx).Take(&user, "id = ?", userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
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

func (ur *userRepository) UpdateUser(ctx context.Context, userId uuid.UUID, user *models.User) error {
	result := ur.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userId).Updates(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user	with ID %s not found", userId)
	}

	return nil
}

func (ur *userRepository) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	result := ur.db.WithContext(ctx).Unscoped().Where("id = ?", userId).Delete(&models.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %s not found", userId)
	}

	return nil
}
