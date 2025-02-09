package models

import "github.com/google/uuid"

type Book struct {
	BaseGormModel
	SecondaryID   string      `json:"secondary_id" gorm:"unique;default:null"`
	CategoryID    uuid.UUID   `json:"category_id" gorm:"not null"`
	SellerID      uuid.UUID   `json:"seller_id" gorm:"not null"`
	Name          string      `json:"name" gorm:"not null;size:255;check:length(name) > 0"`
	Description   string      `json:"description" gorm:"default:null"`
	Price         float64     `json:"price" gorm:"not null;check:price >= 0"`
	OriginalPrice float64     `json:"original_price" gorm:"check:original_price >= 0 AND original_price >= price"`
	RatingAverage float64     `json:"rating_average" gorm:"check:rating_average BETWEEN 0 AND 5"`
	ReviewCount   int         `json:"review_count" gorm:"check:review_count >= 0"`
	PageCount     int         `json:"page_count" gorm:"check:page_count >= 0"`
	Inventory     Inventory   `json:"inventory" gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Authors       []Author    `json:"authors" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;many2many:book_authors"`
	Promotions    []Promotion `json:"promotions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;many2many:promotion_books"`
}

type BookAuthor struct {
	BaseGormModel
	BookID   uuid.UUID `json:"book_id" gorm:"not null"`
	AuthorID uuid.UUID `json:"author_id" gorm:"not null"`
}
