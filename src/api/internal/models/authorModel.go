package models

type Author struct {
	BaseGormModel
	SecondaryID string `json:"secondary_id" gorm:"unique;default:null"`
	Name        string `json:"name"         gorm:"not null"`
	Slug        string `json:"slug"         gorm:"not null"`
	Books       []Book `json:"books"        gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;many2many:book_authors"`
}
