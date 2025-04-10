package models

type Seller struct {
	BaseGormModel
	SecondaryID string `json:"secondary_id" gorm:"unique;default:null"`
	Name        string `json:"name"         gorm:"not null;size:255;check:length(name) > 0"`
	Link        string `json:"link"         gorm:"unique;not null"`
	Logo        string `json:"logo"         gorm:"unique;null"`
	Books       []Book `                    gorm:"foreignKey:SellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
