package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	AdminRole Role = "ADMIN"
	UserRole  Role = "USER"
)

type User struct {
	BaseGormModel
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	LastLogin   time.Time `json:"last_login"`
	IsActive    bool      `json:"is_active"    gorm:"default:true"`
	PhoneNumber string    `json:"phone_number"`
	Role        Role      `json:"role"         gorm:"type:enum('admin', 'user');default:'USER'"`
	Orders      []Order   `json:"orders"       gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	return nil
}
