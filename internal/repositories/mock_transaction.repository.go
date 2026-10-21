package repositories

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"POS/internal/models"
)

type MockTransactionRepository struct {
	Transactions map[uint]*models.Transaction
	Items        map[uint]*models.TransactionItem
	NextID       uint
	NextItemID   uint
	ProductRepo  *MockProductRepository
	StockRepo    *MockStockRepository
}

func NewMockTransactionRepository(prodRepo *MockProductRepository, stockRepo *MockStockRepository) *MockTransactionRepository {
	return &MockTransactionRepository{
		Transactions: make(map[uint]*models.Transaction),
		Items:        make(map[uint]*models.TransactionItem),
		NextID:       1,
		NextItemID:   1,
		ProductRepo:  prodRepo,
		StockRepo:    stockRepo,
	}
}

func (m *MockTransactionRepository) CreateWithItems(_ context.Context, txn *models.Transaction, items []*models.TransactionItem) error {
	// Cek invoice duplikat
	for _, t := range m.Transactions {
		if t.InvoiceNumber == txn.InvoiceNumber && !t.DeletedAt.Valid {
			return errors.New("duplicate key invoice_number")
		}
	}

	// Validasi & kurangi stok
	for _, item := range items {
		product, ok := m.ProductRepo.Products[item.ProductID]
		if !ok || product.DeletedAt.Valid {
			return errors.New("produk tidak ditemukan: " + item.ProductName)
		}
		if product.Stock < item.Quantity {
			return errors.New("stok tidak mencukupi untuk " + product.Name)
		}

		stockBefore := product.Stock
		product.Stock -= item.Quantity

		mutation := &models.StockMutation{
			ProductID:     product.ID,
			Type:          models.StockTypeOut,
			Quantity:      item.Quantity,
			StockBefore:   stockBefore,
			StockAfter:    product.Stock,
			UserID:        txn.CashierID,
			ReferenceType: "transaction",
			Note:          "Penjualan " + txn.InvoiceNumber,
		}
		_ = m.StockRepo.Create(nil, mutation)
	}

	// Simpan transaksi
	txn.ID = m.NextID
	m.NextID++

	for _, item := range items {
		item.ID = m.NextItemID
		m.NextItemID++
		item.TransactionID = txn.ID
		m.Items[item.ID] = item
	}
	txn.Items = items
	m.Transactions[txn.ID] = txn

	// Update reference_id di mutasi
	for _, mut := range m.StockRepo.Mutations {
		if mut.ReferenceType == "transaction" && mut.ReferenceID == 0 && mut.Note == "Penjualan "+txn.InvoiceNumber {
			mut.ReferenceID = txn.ID
		}
	}

	return nil
}

func (m *MockTransactionRepository) Pay(_ context.Context, txn *models.Transaction) error {
	existing, ok := m.Transactions[txn.ID]
	if !ok {
		return errors.New("transaction not found")
	}
	if existing.Status != models.TransactionStatusPending {
		return errors.New("hanya transaksi pending yang bisa dibayar")
	}
	existing.Status = models.TransactionStatusPaid
	existing.PaymentMethod = txn.PaymentMethod
	existing.PaymentAmount = txn.PaymentAmount
	existing.ChangeAmount = txn.ChangeAmount
	existing.PaidAt = txn.PaidAt
	return nil
}

func (m *MockTransactionRepository) VoidWithStockRestore(_ context.Context, txn *models.Transaction, reason string, userID uint) error {
	existing, ok := m.Transactions[txn.ID]
	if !ok {
		return errors.New("transaction not found")
	}

	// Restore stok
	for _, item := range existing.Items {
		product, ok := m.ProductRepo.Products[item.ProductID]
		if !ok {
			continue
		}
		stockBefore := product.Stock
		product.Stock += item.Quantity

		mutation := &models.StockMutation{
			ProductID:     product.ID,
			Type:          models.StockTypeIn,
			Quantity:      item.Quantity,
			StockBefore:   stockBefore,
			StockAfter:    product.Stock,
			UserID:        userID,
			ReferenceID:   existing.ID,
			ReferenceType: "void",
			Note:          "Void " + existing.InvoiceNumber + ": " + reason,
		}
		_ = m.StockRepo.Create(nil, mutation)
	}

	now := time.Now()
	existing.Status = models.TransactionStatusVoid
	existing.VoidedAt = &now
	existing.VoidedBy = &userID
	existing.VoidReason = reason
	return nil
}

func (m *MockTransactionRepository) FindByID(_ context.Context, id uint) (*models.Transaction, error) {
	txn, ok := m.Transactions[id]
	if !ok || txn.DeletedAt.Valid {
		return nil, errors.New("transaction not found")
	}
	return txn, nil
}

func (m *MockTransactionRepository) FindAll(_ context.Context, filter TransactionFilter) ([]*models.Transaction, int64, error) {
	var result []*models.Transaction
	for _, t := range m.Transactions {
		if t.DeletedAt.Valid {
			continue
		}
		if filter.Status != "" && t.Status != filter.Status {
			continue
		}
		if filter.CashierID != nil && t.CashierID != *filter.CashierID {
			continue
		}
		if filter.Search != "" {
			if !strings.Contains(strings.ToLower(t.InvoiceNumber), strings.ToLower(filter.Search)) {
				continue
			}
		}
		if filter.From != nil && t.CreatedAt.Before(*filter.From) {
			continue
		}
		if filter.To != nil && !t.CreatedAt.Before(*filter.To) {
			continue
		}
		result = append(result, t)
	}

	sort.Slice(result, func(i, j int) bool { return result[i].ID > result[j].ID })

	total := int64(len(result))

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit
	if offset >= len(result) {
		return []*models.Transaction{}, total, nil
	}
	end := offset + filter.Limit
	if end > len(result) {
		end = len(result)
	}

	return result[offset:end], total, nil
}

func (m *MockTransactionRepository) CountToday(_ context.Context) (int64, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)

	var count int64
	for _, t := range m.Transactions {
		if t.CreatedAt.After(start) && t.CreatedAt.Before(end) {
			count++
		}
	}
	return count, nil
}
