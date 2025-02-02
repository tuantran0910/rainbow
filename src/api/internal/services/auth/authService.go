package auth

import (
	"context"
	"fmt"

	model "github.com/tuantran0910/rainbow/internal/models/user"
	repository "github.com/tuantran0910/rainbow/internal/repositories/user"
	"github.com/tuantran0910/rainbow/pkg/utils/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAuthService interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, userRequest model.UserRequest) error
}

type authService struct {
	db             *gorm.DB
	userRepository repository.IUserRepository
}

func NewAuthService(db *gorm.DB, userRepository repository.IUserRepository) IAuthService {
	return &authService{
		db:             db,
		userRepository: userRepository,
	}
}

func (as *authService) withTx(ctx context.Context, fn func(context.Context, repository.IUserRepository) error) error {
	return as.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userRepository := as.userRepository.WithTX(tx)
		return fn(ctx, userRepository)
	})
}

func (as *authService) Login(ctx context.Context, email string, password string) (string, error) {
	// Retrieve the user by email
	user, err := as.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// Compare the provided password and the password in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}

	return token, nil
}

func (as *authService) Register(ctx context.Context, userRequest model.UserRequest) error {
	// Check if the email already exists
	existingUser, err := as.userRepository.GetUserByEmail(ctx, *userRequest.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existingUser != nil {
		return fmt.Errorf("email already exists")
	}

	return as.withTx(ctx, func(ctx context.Context, userRepository repository.IUserRepository) error {
		// Define a new user
		user := &model.User{
			Email:       *userRequest.Email,
			Password:    *userRequest.Password,
			FirstName:   *userRequest.FirstName,
			LastName:    *userRequest.LastName,
			PhoneNumber: *userRequest.PhoneNumber,
		}

		// Create a new user
		if err := userRepository.CreateUser(ctx, user); err != nil {
			return err
		}

		return nil
	})
}
