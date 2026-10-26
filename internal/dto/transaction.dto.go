package dto

import "time"

type CreateTransactionItemRequest struct {
	ProductID uint `json:"product_id" binding:"required" example:"1"`
	Quantity  int  `json:"quantity" binding:"required,gt=0" example:"2"`
}

type CreateTransactionRequest struct {
	Items         []CreateTransactionItemRequest `json:"items" binding:"required,min=1,dive"`
	DiscountType  string                         `json:"discount_type" binding:"omitempty,oneof=percentage fixed" example:"percentage"`
	DiscountValue float64                        `json:"discount_value" binding:"gte=0" example:"10"`
	PaymentMethod string                         `json:"payment_method" binding:"omitempty,oneof=cash qris card transfer" example:"cash"`
	PaymentAmount float64                        `json:"payment_amount" binding:"gte=0" example:"20000"`
	Note          string                         `json:"note" binding:"max=500" example:"Pelanggan reguler"`
}

type PayTransactionRequest struct {
	PaymentMethod string  `json:"payment_method" binding:"required,oneof=cash qris card transfer" example:"cash"`
	PaymentAmount float64 `json:"payment_amount" binding:"required,gt=0" example:"20000"`
}

type VoidTransactionRequest struct {
	Reason string `json:"reason" binding:"required,max=255" example:"Customer batal"`
}

type TransactionItemResponse struct {
	ProductID   uint    `json:"product_id"`
	ProductSKU  string  `json:"product_sku"`
	ProductName string  `json:"product_name"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

type TransactionResponse struct {
	ID             uint                      `json:"id"`
	InvoiceNumber  string                    `json:"invoice_number"`
	CashierID      uint                      `json:"cashier_id"`
	Status         string                    `json:"status"`
	Subtotal       float64                   `json:"subtotal"`
	DiscountType   string                    `json:"discount_type,omitempty"`
	DiscountValue  float64                   `json:"discount_value,omitempty"`
	DiscountAmount float64                   `json:"discount_amount"`
	TaxRate        float64                   `json:"tax_rate"`
	TaxAmount      float64                   `json:"tax_amount"`
	Total          float64                   `json:"total"`
	PaymentMethod  string                    `json:"payment_method,omitempty"`
	PaymentAmount  float64                   `json:"payment_amount,omitempty"`
	ChangeAmount   float64                   `json:"change_amount,omitempty"`
	PaidAt         *time.Time                `json:"paid_at,omitempty"`
	VoidedAt       *time.Time                `json:"voided_at,omitempty"`
	VoidedBy       *uint                     `json:"voided_by,omitempty"`
	VoidReason     string                    `json:"void_reason,omitempty"`
	Note           string                    `json:"note,omitempty"`
	Items          []TransactionItemResponse `json:"items"`
	CreatedAt      time.Time                 `json:"created_at"`
}
