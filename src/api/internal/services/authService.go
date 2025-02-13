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

func (as *authService) withTX(
	ctx context.Context,
	fn func(context.Context, repositories.IUserRepository) error,
) error {
	return as.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userRepository := as.userRepository.WithTX(tx)
		return fn(ctx, userRepository)
	})
}

func (as *authService) Login(
	ctx context.Context,
	loginUserRequest dtos.LoginUserRequest,
) (string, error) {
	user, err := as.userRepository.GetUserByEmail(ctx, loginUserRequest.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUserRequest.Password)); err != nil {
		return "", fmt.Errorf("invalid credentials: %w", err)
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return token, nil
}

func (as *authService) Register(
	ctx context.Context,
	registerUserRequest dtos.RegisterUserRequest,
) error {
	return as.withTX(
		ctx,
		func(ctx context.Context, userRepository repositories.IUserRepository) error {
			existingUserWithEmail, err := as.userRepository.GetUserByEmail(
				ctx,
				registerUserRequest.Email,
			)
			if err != nil {
				if existingUserWithEmail != nil {
					return fmt.Errorf("email already exists")
				}
				return err
			}

			existingUserWithPhoneNumber, err := as.userRepository.GetUserByPhoneNumber(
				ctx,
				registerUserRequest.PhoneNumber,
			)
			if err != nil {
				if existingUserWithPhoneNumber != nil {
					return fmt.Errorf("phone number already exists")
				}
				return err
			}

			user := &models.User{
				Email:       registerUserRequest.Email,
				Password:    registerUserRequest.Password,
				FirstName:   registerUserRequest.FirstName,
				LastName:    registerUserRequest.LastName,
				PhoneNumber: registerUserRequest.PhoneNumber,
				Role:        models.Role(registerUserRequest.Role),
			}
			return userRepository.CreateUser(ctx, user)
		},
	)
}
