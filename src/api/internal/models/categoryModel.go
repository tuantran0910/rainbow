package models

type Category struct {
	BaseGormModel
	Name  string `json:"name" gorm:"not null;size:255;check:length(name) > 0"`
	Slug  string `json:"slug" gorm:"unique;not null;size:255;check:length(slug) > 0"`
	Books []Book `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
