package dto

import "time"

type CreateStockMutationRequest struct {
	ProductID uint   `json:"product_id" binding:"required" example:"1"`
	Type      string `json:"type" binding:"required,oneof=in out adjustment" example:"in"`
	Quantity  int    `json:"quantity" binding:"required,gte=0" example:"50"`
	Note      string `json:"note" binding:"max=255" example:"Restock dari supplier"`
}

type StockMutationResponse struct {
	ID          uint             `json:"id"`
	ProductID   uint             `json:"product_id"`
	Product     *ProductResponse `json:"product,omitempty"`
	Type        string           `json:"type"`
	Quantity    int              `json:"quantity"`
	StockBefore int              `json:"stock_before"`
	StockAfter  int              `json:"stock_after"`
	Note        string           `json:"note"`
	UserID      uint             `json:"user_id"`
	CreatedAt   time.Time        `json:"created_at"`
}
