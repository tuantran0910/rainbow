package models

type Author struct {
	BaseGormModel
	Name string `json:"name" gorm:"not null"`
	Slug string `json:"slug" gorm:"unique;not null"`
}
