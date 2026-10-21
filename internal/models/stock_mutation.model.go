package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	StockTypeIn         = "in"
	StockTypeOut        = "out"
	StockTypeAdjustment = "adjustment"
)

type StockMutation struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ProductID     uint           `gorm:"not null;index" json:"product_id"`
	Product       *Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Type          string         `gorm:"size:20;not null;index" json:"type"`
	Quantity      int            `gorm:"not null" json:"quantity"`
	StockBefore   int            `gorm:"not null" json:"stock_before"`
	StockAfter    int            `gorm:"not null" json:"stock_after"`
	Note          string         `gorm:"size:255" json:"note"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	ReferenceID   uint           `gorm:"index" json:"reference_id,omitempty"`
	ReferenceType string         `gorm:"size:30;index" json:"reference_type,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
