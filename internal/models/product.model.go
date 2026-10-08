package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	SKU         string         `gorm:"size:50;not null" json:"sku"`
	Name        string         `gorm:"size:150;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Price       float64        `gorm:"not null;default:0" json:"price"`
	Cost        float64        `gorm:"not null;default:0" json:"cost"`
	Stock       int            `gorm:"not null;default:0" json:"stock"`
	CategoryID  uint           `gorm:"not null;index" json:"category_id"`
	Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	IsActive    bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
