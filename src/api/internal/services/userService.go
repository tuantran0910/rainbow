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

type IUserService interface {
	GetUsers(ctx context.Context, page, limit int) ([]*models.User, *pagination.Pagination, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, userId uuid.UUID, userRequest dtos.UpdateUserRequest, currentUserId uuid.UUID) error
	DeleteUser(ctx context.Context, userId uuid.UUID, currentUserId uuid.UUID) error
}

type userService struct {
	db             *gorm.DB
	userRepository repositories.IUserRepository
}

func NewUserService(db *gorm.DB, userRepository repositories.IUserRepository) IUserService {
	return &userService{
		db:             db,
		userRepository: userRepository,
	}
}

func (us *userService) withTx(ctx context.Context, fn func(context.Context, repositories.IUserRepository) error) error {
	return us.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userRepository := us.userRepository.WithTX(tx)
		return fn(ctx, userRepository)
	})
}

func (us *userService) GetUsers(ctx context.Context, page, limit int) ([]*models.User, *pagination.Pagination, error) {
	return us.userRepository.GetUsers(ctx, page, limit)
}

func (us *userService) GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	return us.userRepository.GetUserById(ctx, userId)
}

//gocyclo:ignore
func (us *userService) UpdateUser(ctx context.Context, userId uuid.UUID, userRequest dtos.UpdateUserRequest, currentUserId uuid.UUID) error {
	return us.withTx(ctx, func(ctx context.Context, userRepository repositories.IUserRepository) error {
		// Get user by id
		user, err := userRepository.GetUserById(ctx, userId)
		if err != nil {
			return err
		}

		// Get current user by id
		currentUser, err := userRepository.GetUserById(ctx, currentUserId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if currentUser == nil {
			return fmt.Errorf("current user not found")
		}

		if currentUser.ID != userId && string(currentUser.Role) != string(models.AdminRole) {
			return fmt.Errorf("only admin can update other user's information")
		}

		// Apply the updates
		if userRequest.Email != nil && *userRequest.Email != user.Email {
			// Check if the email is already taken
			userByEmail, err := userRepository.GetUserByEmail(ctx, *userRequest.Email)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			if userByEmail != nil && userByEmail.ID != userId {
				return fmt.Errorf("email %s is already taken", *userRequest.Email)
			}

			user.Email = *userRequest.Email
		}

		if userRequest.PhoneNumber != nil && *userRequest.PhoneNumber != user.PhoneNumber {
			// Check if the phone number is already taken
			userByPhoneNumber, err := userRepository.GetUserByPhoneNumber(ctx, *userRequest.PhoneNumber)
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}

			if userByPhoneNumber != nil && userByPhoneNumber.ID != userId {
				return fmt.Errorf("phone number %s is already taken", *userRequest.PhoneNumber)
			}

			user.PhoneNumber = *userRequest.PhoneNumber
		}

		if userRequest.FirstName != nil && *userRequest.FirstName != user.FirstName {
			user.FirstName = *userRequest.FirstName
		}

		if userRequest.LastName != nil && *userRequest.LastName != user.LastName {
			user.LastName = *userRequest.LastName
		}

		if userRequest.IsActive != nil && *userRequest.IsActive != user.IsActive {
			user.IsActive = *userRequest.IsActive
		}

		if userRequest.Role != nil && *userRequest.Role != string(user.Role) {
			// Check if the current user is admin
			if string(currentUser.Role) == string(models.AdminRole) {
				user.Role = models.Role(*userRequest.Role)
			}

			return fmt.Errorf("only admin can update user role")
		}

		if err := userRepository.UpdateUser(ctx, userId, user); err != nil {
			return err
		}

		return nil
	})
}

func (us *userService) DeleteUser(ctx context.Context, userId uuid.UUID, currentUserId uuid.UUID) error {
	return us.withTx(ctx, func(ctx context.Context, userRepository repositories.IUserRepository) error {
		// Get the current user by id
		currentUser, err := userRepository.GetUserById(ctx, currentUserId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if currentUser == nil {
			return fmt.Errorf("current user not found")
		}

		// Check if the current user is admin
		if string(currentUser.Role) != string(models.AdminRole) {
			return fmt.Errorf("only admin can delete user")
		}

		// Delete a user
		return userRepository.DeleteUser(ctx, userId)
	})
}
