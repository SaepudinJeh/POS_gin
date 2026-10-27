package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"POS/internal/models"
	"POS/internal/repositories"
)

type TransactionItemInput struct {
	ProductID uint
	Quantity  int
}

type TransactionInput struct {
	Items         []TransactionItemInput
	DiscountType  string
	DiscountValue float64
	PaymentMethod string
	PaymentAmount float64
	Note          string
	CashierID     uint
}

type TransactionService struct {
	txnRepo     repositories.TransactionRepository
	productRepo repositories.ProductRepository
	settingSvc  *SettingService
}

func NewTransactionService(
	txnRepo repositories.TransactionRepository,
	productRepo repositories.ProductRepository,
	settingSvc *SettingService,
) *TransactionService {
	return &TransactionService{
		txnRepo:     txnRepo,
		productRepo: productRepo,
		settingSvc:  settingSvc,
	}
}

func (s *TransactionService) GetAll(ctx context.Context, filter repositories.TransactionFilter) ([]*models.Transaction, int64, error) {
	return s.txnRepo.FindAll(ctx, filter)
}

func (s *TransactionService) GetByID(ctx context.Context, id uint) (*models.Transaction, error) {
	return s.txnRepo.FindByID(ctx, id)
}

func (s *TransactionService) Create(ctx context.Context, input TransactionInput) (*models.Transaction, error) {
	if len(input.Items) == 0 {
		return nil, errors.New("transaksi harus punya minimal 1 item")
	}

	// Cek duplikat produk
	seen := make(map[uint]bool)
	for _, item := range input.Items {
		if seen[item.ProductID] {
			return nil, errors.New("produk duplikat dalam transaksi")
		}
		seen[item.ProductID] = true
	}

	// Load produk + validasi stok
	var items []*models.TransactionItem
	var subtotal float64

	for _, in := range input.Items {
		product, err := s.productRepo.FindByID(ctx, in.ProductID)
		if err != nil {
			return nil, errors.New("produk tidak ditemukan")
		}
		if !product.IsActive {
			return nil, errors.New("produk tidak aktif: " + product.Name)
		}
		if product.Stock < in.Quantity {
			return nil, errors.New("stok tidak mencukupi untuk " + product.Name)
		}

		itemSubtotal := product.Price * float64(in.Quantity)
		subtotal += itemSubtotal

		items = append(items, &models.TransactionItem{
			ProductID:   product.ID,
			ProductSKU:  product.SKU,
			ProductName: product.Name,
			UnitPrice:   product.Price,
			UnitCost:    product.Cost,
			Quantity:    in.Quantity,
			Subtotal:    itemSubtotal,
		})
	}

	// Hitung diskon
	discountAmount := 0.0
	if input.DiscountType != "" && input.DiscountValue > 0 {
		switch input.DiscountType {
		case models.DiscountTypePercentage:
			if input.DiscountValue > 100 {
				return nil, errors.New("diskon persentase tidak boleh lebih dari 100")
			}
			discountAmount = subtotal * input.DiscountValue / 100
		case models.DiscountTypeFixed:
			discountAmount = input.DiscountValue
		default:
			return nil, errors.New("tipe diskon tidak valid")
		}
		if discountAmount > subtotal {
			discountAmount = subtotal
		}
	}

	// Hitung pajak dari setting
	taxRate := s.settingSvc.GetTaxRate(ctx)
	taxAmount := (subtotal - discountAmount) * taxRate / 100

	total := subtotal - discountAmount + taxAmount

	// Bangun transaksi
	now := time.Now()
	txn := &models.Transaction{
		CashierID:      input.CashierID,
		Subtotal:       subtotal,
		DiscountType:   input.DiscountType,
		DiscountValue:  input.DiscountValue,
		DiscountAmount: discountAmount,
		TaxRate:        taxRate,
		TaxAmount:      taxAmount,
		Total:          total,
		Note:           input.Note,
		Status:         models.TransactionStatusPending,
		CreatedAt:      now,
	}

	// Kalau payment langsung diisi → paid
	if input.PaymentMethod != "" && input.PaymentAmount > 0 {
		if input.PaymentAmount < total {
			return nil, errors.New("jumlah pembayaran kurang dari total")
		}
		txn.Status = models.TransactionStatusPaid
		txn.PaymentMethod = input.PaymentMethod
		txn.PaymentAmount = input.PaymentAmount
		txn.ChangeAmount = input.PaymentAmount - total
		txn.PaidAt = &now
	}

	// Generate invoice + retry kalau duplikat
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		invoice, err := s.generateInvoice(ctx)
		if err != nil {
			return nil, err
		}
		txn.InvoiceNumber = invoice

		lastErr = s.txnRepo.CreateWithItems(ctx, txn, items)
		if lastErr == nil {
			return s.txnRepo.FindByID(ctx, txn.ID)
		}
		if !repositories.IsDuplicateKey(lastErr) {
			return nil, lastErr
		}
	}
	return nil, fmt.Errorf("gagal generate invoice unik: %w", lastErr)
}

func (s *TransactionService) Pay(ctx context.Context, id uint, method string, amount float64) (*models.Transaction, error) {
	txn, err := s.txnRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}
	if txn.Status == models.TransactionStatusVoid {
		return nil, errors.New("transaksi sudah void")
	}
	if txn.Status == models.TransactionStatusPaid {
		return nil, errors.New("transaksi sudah dibayar")
	}
	if amount < txn.Total {
		return nil, errors.New("jumlah pembayaran kurang dari total")
	}

	now := time.Now()
	txn.PaymentMethod = method
	txn.PaymentAmount = amount
	txn.ChangeAmount = amount - txn.Total
	txn.PaidAt = &now

	if err := s.txnRepo.Pay(ctx, txn); err != nil {
		return nil, err
	}

	return s.txnRepo.FindByID(ctx, id)
}

func (s *TransactionService) Void(ctx context.Context, id uint, reason string, userID uint) (*models.Transaction, error) {
	txn, err := s.txnRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}
	if txn.Status == models.TransactionStatusVoid {
		return nil, errors.New("transaksi sudah void")
	}

	if err := s.txnRepo.VoidWithStockRestore(ctx, txn, reason, userID); err != nil {
		return nil, err
	}

	return s.txnRepo.FindByID(ctx, id)
}

func (s *TransactionService) generateInvoice(ctx context.Context) (string, error) {
	count, err := s.txnRepo.CountToday(ctx)
	if err != nil {
		return "", err
	}
	prefix := s.settingSvc.GetString(ctx, SettingKeyInvoicePrefix, "INV")
	date := time.Now().Format("20060102")
	return fmt.Sprintf("%s-%s-%04d", prefix, date, count+1), nil
}
