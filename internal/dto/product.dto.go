package dto

import "time"

type CreateProductRequest struct {
	SKU         string  `json:"sku" binding:"required,min=2,max=50,notspace"`
	Name        string  `json:"name" binding:"required,min=2,max=150"`
	Description string  `json:"description" binding:"max=500"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	Cost        float64 `json:"cost" binding:"gte=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

type UpdateProductRequest struct {
	SKU         string  `json:"sku" binding:"required,min=2,max=50,notspace"`
	Name        string  `json:"name" binding:"required,min=2,max=150"`
	Description string  `json:"description" binding:"max=500"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	Cost        float64 `json:"cost" binding:"gte=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	CategoryID  uint    `json:"category_id" binding:"required"`
	IsActive    *bool   `json:"is_active"`
}

type ProductResponse struct {
	ID          uint              `json:"id"`
	SKU         string            `json:"sku"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Cost        float64           `json:"cost"`
	Stock       int               `json:"stock"`
	IsActive    bool              `json:"is_active"`
	Category    *CategoryResponse `json:"category,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
