package repositories

import (
	"context"
	"errors"
	"strings"
	"time"

	"POS/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionFilter struct {
	Status    string
	CashierID *uint
	Search    string
	From      *time.Time
	To        *time.Time
	Page      int
	Limit     int
}

type TransactionRepository interface {
	CreateWithItems(ctx context.Context, txn *models.Transaction, items []*models.TransactionItem) error
	Pay(ctx context.Context, txn *models.Transaction) error
	VoidWithStockRestore(ctx context.Context, txn *models.Transaction, reason string, userID uint) error
	FindByID(ctx context.Context, id uint) (*models.Transaction, error)
	FindAll(ctx context.Context, filter TransactionFilter) ([]*models.Transaction, int64, error)
	CountToday(ctx context.Context) (int64, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateWithItems(ctx context.Context, txn *models.Transaction, items []*models.TransactionItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock semua produk terkait
		productIDs := make([]uint, 0, len(items))
		for _, item := range items {
			productIDs = append(productIDs, item.ProductID)
		}
		// Sort supaya tidak deadlock
		var products []*models.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id IN ?", productIDs).
			Order("id ASC").
			Find(&products).Error; err != nil {
			return err
		}

		productMap := make(map[uint]*models.Product)
		for _, p := range products {
			productMap[p.ID] = p
		}

		// Kurangi stok + catat mutasi
		for _, item := range items {
			product, ok := productMap[item.ProductID]
			if !ok {
				return errors.New("produk tidak ditemukan: " + item.ProductName)
			}
			if product.Stock < item.Quantity {
				return errors.New("stok tidak mencukupi untuk " + product.Name)
			}

			stockBefore := product.Stock
			stockAfter := stockBefore - item.Quantity

			if err := tx.Model(product).Update("stock", stockAfter).Error; err != nil {
				return err
			}

			mutation := &models.StockMutation{
				ProductID:     product.ID,
				Type:          models.StockTypeOut,
				Quantity:      item.Quantity,
				StockBefore:   stockBefore,
				StockAfter:    stockAfter,
				UserID:        txn.CashierID,
				ReferenceType: "transaction",
				Note:          "Penjualan " + txn.InvoiceNumber,
			}
			if err := tx.Create(mutation).Error; err != nil {
				return err
			}
		}

		// Buat transaksi
		if err := tx.Create(txn).Error; err != nil {
			return err
		}

		// Buat items
		for _, item := range items {
			item.TransactionID = txn.ID
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}

		// Update ReferenceID di mutasi
		if err := tx.Model(&models.StockMutation{}).
			Where("reference_type = ? AND reference_id = 0 AND note = ?",
				"transaction", "Penjualan "+txn.InvoiceNumber).
			Update("reference_id", txn.ID).Error; err != nil {
			return err
		}

		txn.Items = items
		return nil
	})
}

func (r *transactionRepository) Pay(ctx context.Context, txn *models.Transaction) error {
	return r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("id = ? AND status = ?", txn.ID, models.TransactionStatusPending).
		Updates(map[string]interface{}{
			"status":         models.TransactionStatusPaid,
			"payment_method": txn.PaymentMethod,
			"payment_amount": txn.PaymentAmount,
			"change_amount":  txn.ChangeAmount,
			"paid_at":        txn.PaidAt,
		}).Error
}

func (r *transactionRepository) VoidWithStockRestore(ctx context.Context, txn *models.Transaction, reason string, userID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Restore stok + catat mutasi
		for _, item := range txn.Items {
			var product models.Product
			if err := tx.Unscoped().
				Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&product, item.ProductID).Error; err != nil {
				return err
			}

			stockBefore := product.Stock
			stockAfter := stockBefore + item.Quantity

			if err := tx.Unscoped().Model(&product).Update("stock", stockAfter).Error; err != nil {
				return err
			}

			mutation := &models.StockMutation{
				ProductID:     product.ID,
				Type:          models.StockTypeIn,
				Quantity:      item.Quantity,
				StockBefore:   stockBefore,
				StockAfter:    stockAfter,
				UserID:        userID,
				ReferenceID:   txn.ID,
				ReferenceType: "void",
				Note:          "Void " + txn.InvoiceNumber + ": " + reason,
			}
			if err := tx.Create(mutation).Error; err != nil {
				return err
			}
		}

		now := time.Now()
		return tx.Model(&models.Transaction{}).
			Where("id = ?", txn.ID).
			Updates(map[string]interface{}{
				"status":      models.TransactionStatusVoid,
				"voided_at":   &now,
				"voided_by":   userID,
				"void_reason": reason,
			}).Error
	})
}

func (r *transactionRepository) FindByID(ctx context.Context, id uint) (*models.Transaction, error) {
	var txn models.Transaction
	if err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Cashier").
		First(&txn, id).Error; err != nil {
		return nil, err
	}
	return &txn, nil
}

func (r *transactionRepository) FindAll(ctx context.Context, filter TransactionFilter) ([]*models.Transaction, int64, error) {
	var txns []*models.Transaction
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Transaction{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.CashierID != nil {
		query = query.Where("cashier_id = ?", *filter.CashierID)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("invoice_number ILIKE ?", search)
	}
	if filter.From != nil {
		query = query.Where("created_at >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("created_at < ?", *filter.To)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	if err := query.
		Preload("Items").
		Order("id DESC").
		Offset(offset).
		Limit(filter.Limit).
		Find(&txns).Error; err != nil {
		return nil, 0, err
	}

	return txns, total, nil
}

func (r *transactionRepository) CountToday(ctx context.Context) (int64, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)

	var count int64
	if err := r.db.WithContext(ctx).Unscoped().
		Model(&models.Transaction{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// IsDuplicateKey helper untuk retry invoice
func IsDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "SQLSTATE 23505")
}
