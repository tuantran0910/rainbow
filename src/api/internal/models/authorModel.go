package models

type Author struct {
	BaseGormModel
	Name  string `json:"name" gorm:"not null"`
	Slug  string `json:"slug" gorm:"unique;not null"`
	Books []Book `json:"books" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;many2many:book_authors"`
}
