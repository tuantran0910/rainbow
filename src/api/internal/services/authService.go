package services

import (
	"context"
	"fmt"

	"github.com/tuantran0910/rainbow/internal/dtos"
	"github.com/tuantran0910/rainbow/internal/models"
	"github.com/tuantran0910/rainbow/internal/repositories"
	"github.com/tuantran0910/rainbow/pkg/utils/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAuthService interface {
	Login(ctx context.Context, loginUserRequest dtos.LoginUserRequest) (string, error)
	Register(ctx context.Context, registerUserRequest dtos.RegisterUserRequest) error
}

type authService struct {
	db             *gorm.DB
	userRepository repositories.IUserRepository
}

func NewAuthService(db *gorm.DB, userRepository repositories.IUserRepository) IAuthService {
	return &authService{
		db:             db,
		userRepository: userRepository,
	}
}

func (as *authService) withTx(ctx context.Context, fn func(context.Context, repositories.IUserRepository) error) error {
	return as.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userRepository := as.userRepository.WithTX(tx)
		return fn(ctx, userRepository)
	})
}

func (as *authService) Login(ctx context.Context, loginUserRequest dtos.LoginUserRequest) (string, error) {
	// Retrieve the user by email
	user, err := as.userRepository.GetUserByEmail(ctx, loginUserRequest.Email)
	if err != nil {
		return "", err
	}

	// Compare the provided password and the password in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUserRequest.Password)); err != nil {
		return "", fmt.Errorf("invalid credentials: %w", err)
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (as *authService) Register(ctx context.Context, registerUserRequest dtos.RegisterUserRequest) error {
	var existingUser *models.User

	// Check if the email already exists
	existingUser, err := as.userRepository.GetUserByEmail(ctx, registerUserRequest.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existingUser != nil {
		return fmt.Errorf("email already exists")
	}

	// Check if the phone number already exists
	existingUser, err = as.userRepository.GetUserByPhoneNumber(ctx, registerUserRequest.PhoneNumber)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existingUser != nil {
		return fmt.Errorf("phone number already exists")
	}

	return as.withTx(ctx, func(ctx context.Context, userRepository repositories.IUserRepository) error {
		// Define a new user
		user := &models.User{
			Email:       registerUserRequest.Email,
			Password:    registerUserRequest.Password,
			FirstName:   registerUserRequest.FirstName,
			LastName:    registerUserRequest.LastName,
			PhoneNumber: registerUserRequest.PhoneNumber,
			Role:        models.Role(registerUserRequest.Role),
		}

		// Create a new user
		if err := userRepository.CreateUser(ctx, user); err != nil {
			return err
		}

		return nil
	})
}
