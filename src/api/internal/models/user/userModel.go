package user

import (
	"time"

	"github.com/tuantran0910/rainbow/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	models.BaseGormModel
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	LastLogin   time.Time `json:"last_login"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	PhoneNumber string    `json:"phone_number"`
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

type UserRequest struct {
	Email       *string    `json:"email"`
	Password    *string    `json:"password"`
	FirstName   *string    `json:"first_name"`
	LastName    *string    `json:"last_name"`
	LastLogin   *time.Time `json:"last_login" gorm:"default:null"`
	IsActive    *bool      `json:"is_active" gorm:"default:true"`
	PhoneNumber *string    `json:"phone_number"`
}
