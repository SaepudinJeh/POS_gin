package repositories

import (
	"context"
	"errors"
	"sort"
	"time"

	"POS/internal/models"

	"gorm.io/gorm"
)

type MockStockRepository struct {
	Mutations   map[uint]*models.StockMutation
	NextID      uint
	ProductRepo *MockProductRepository
}

func NewMockStockRepository(productRepo *MockProductRepository) *MockStockRepository {
	return &MockStockRepository{
		Mutations:   make(map[uint]*models.StockMutation),
		NextID:      1,
		ProductRepo: productRepo,
	}
}

func isActiveMutation(m *models.StockMutation) bool {
	return !m.DeletedAt.Valid
}

func (m *MockStockRepository) preloadProduct(mut *models.StockMutation) {
	if m.ProductRepo == nil {
		return
	}
	if p, err := m.ProductRepo.FindByID(mut.ProductID); err == nil {
		mut.Product = p
	}
}

func (m *MockStockRepository) Create(_ context.Context, mutation *models.StockMutation) error {
	mutation.ID = m.NextID
	m.NextID++
	m.Mutations[mutation.ID] = mutation
	return nil
}

func (m *MockStockRepository) FindAll(_ context.Context, filter StockFilter) ([]*models.StockMutation, int64, error) {
	var result []*models.StockMutation

	for _, mut := range m.Mutations {
		if !isActiveMutation(mut) {
			continue
		}
		if filter.ProductID != nil && mut.ProductID != *filter.ProductID {
			continue
		}
		if filter.Type != "" && mut.Type != filter.Type {
			continue
		}
		m.preloadProduct(mut)
		result = append(result, mut)
	}

	// Sort by ID DESC
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID > result[j].ID
	})

	total := int64(len(result))

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit
	if offset >= len(result) {
		return []*models.StockMutation{}, total, nil
	}
	end := offset + filter.Limit
	if end > len(result) {
		end = len(result)
	}

	return result[offset:end], total, nil
}

func (m *MockStockRepository) FindByID(_ context.Context, id uint) (*models.StockMutation, error) {
	mut, ok := m.Mutations[id]
	if !ok || !isActiveMutation(mut) {
		return nil, errors.New("stock mutation not found")
	}
	m.preloadProduct(mut)
	return mut, nil
}

func (m *MockStockRepository) FindByProductID(_ context.Context, productID uint) ([]*models.StockMutation, error) {
	var result []*models.StockMutation
	for _, mut := range m.Mutations {
		if isActiveMutation(mut) && mut.ProductID == productID {
			m.preloadProduct(mut)
			result = append(result, mut)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID > result[j].ID
	})
	return result, nil
}

// pastikan tidak ada import gorm yang tidak terpakai
var _ = gorm.DeletedAt{}
var _ = time.Now
