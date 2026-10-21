package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	TransactionStatusPending = "pending"
	TransactionStatusPaid    = "paid"
	TransactionStatusVoid    = "void"

	DiscountTypePercentage = "percentage"
	DiscountTypeFixed      = "fixed"

	PaymentMethodCash     = "cash"
	PaymentMethodQRIS     = "qris"
	PaymentMethodCard     = "card"
	PaymentMethodTransfer = "transfer"
)

type Transaction struct {
	ID             uint               `gorm:"primaryKey" json:"id"`
	InvoiceNumber  string             `gorm:"size:30;not null;index" json:"invoice_number"`
	CashierID      uint               `gorm:"not null;index" json:"cashier_id"`
	Cashier        *User              `gorm:"foreignKey:CashierID" json:"cashier,omitempty"`
	Status         string             `gorm:"size:20;not null;index;default:'pending'" json:"status"`
	Subtotal       float64            `gorm:"not null;default:0" json:"subtotal"`
	DiscountType   string             `gorm:"size:20" json:"discount_type,omitempty"`
	DiscountValue  float64            `gorm:"default:0" json:"discount_value,omitempty"`
	DiscountAmount float64            `gorm:"default:0" json:"discount_amount,omitempty"`
	TaxRate        float64            `gorm:"default:0" json:"tax_rate"`
	TaxAmount      float64            `gorm:"default:0" json:"tax_amount"`
	Total          float64            `gorm:"not null;default:0" json:"total"`
	PaymentMethod  string             `gorm:"size:20" json:"payment_method,omitempty"`
	PaymentAmount  float64            `gorm:"default:0" json:"payment_amount,omitempty"`
	ChangeAmount   float64            `gorm:"default:0" json:"change_amount,omitempty"`
	PaidAt         *time.Time         `json:"paid_at,omitempty"`
	VoidedAt       *time.Time         `json:"voided_at,omitempty"`
	VoidedBy       *uint              `json:"voided_by,omitempty"`
	VoidReason     string             `gorm:"size:255" json:"void_reason,omitempty"`
	Note           string             `gorm:"size:500" json:"note,omitempty"`
	Items          []*TransactionItem `gorm:"foreignKey:TransactionID" json:"items,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	DeletedAt      gorm.DeletedAt     `gorm:"index" json:"-"`
}

type TransactionItem struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionID uint      `gorm:"not null;index" json:"transaction_id"`
	ProductID     uint      `gorm:"not null;index" json:"product_id"`
	ProductSKU    string    `gorm:"size:50;not null" json:"product_sku"`
	ProductName   string    `gorm:"size:150;not null" json:"product_name"`
	UnitPrice     float64   `gorm:"not null" json:"unit_price"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	Subtotal      float64   `gorm:"not null" json:"subtotal"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
