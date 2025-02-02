package user

import (
	"context"
	"errors"
	"fmt"

	model "github.com/tuantran0910/rainbow/internal/models/user"
	"gorm.io/gorm"
)

type IUserRepository interface {
	WithTX(tx *gorm.DB) IUserRepository
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
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

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := ur.db.WithContext(ctx).Take(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

func (ur *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	if err := ur.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
