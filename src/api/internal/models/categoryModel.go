package models

type Category struct {
	BaseGormModel
	SecondaryID string `json:"secondary_id" gorm:"unique;default:null"`
	Name        string `json:"name"         gorm:"not null;size:255;check:length(name) > 0"`
	Slug        string `json:"slug"         gorm:"not null;size:255;check:length(slug) > 0"`
	Books       []Book `                    gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
