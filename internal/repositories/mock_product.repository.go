package repositories

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"POS/internal/models"

	"gorm.io/gorm"
)

type MockProductRepository struct {
	Products     map[uint]*models.Product
	NextID       uint
	CategoryRepo *MockCategoryRepository
}

func NewMockProductRepository(categoryRepo *MockCategoryRepository) *MockProductRepository {
	return &MockProductRepository{
		Products:     make(map[uint]*models.Product),
		NextID:       1,
		CategoryRepo: categoryRepo,
	}
}

func isActiveProduct(p *models.Product) bool {
	return !p.DeletedAt.Valid
}

func (m *MockProductRepository) preloadCategory(ctx context.Context, p *models.Product) {
	if m.CategoryRepo == nil {
		return
	}
	if cat, err := m.CategoryRepo.FindByID(ctx, p.CategoryID); err == nil {
		p.Category = cat
	}
}

func (m *MockProductRepository) FindAll(ctx context.Context, filter ProductFilter) ([]*models.Product, int64, error) {
	var result []*models.Product

	for _, p := range m.Products {
		if !isActiveProduct(p) {
			continue
		}
		if filter.CategoryID != nil && p.CategoryID != *filter.CategoryID {
			continue
		}
		if filter.IsActive != nil && p.IsActive != *filter.IsActive {
			continue
		}
		if filter.Search != "" {
			search := strings.ToLower(filter.Search)
			if !strings.Contains(strings.ToLower(p.Name), search) &&
				!strings.Contains(strings.ToLower(p.SKU), search) {
				continue
			}
		}
		m.preloadCategory(ctx, p)
		result = append(result, p)
	}

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
		return []*models.Product{}, total, nil
	}
	end := offset + filter.Limit
	if end > len(result) {
		end = len(result)
	}

	return result[offset:end], total, nil
}

func (m *MockProductRepository) FindByID(ctx context.Context, id uint) (*models.Product, error) {
	p, ok := m.Products[id]
	if !ok || !isActiveProduct(p) {
		return nil, errors.New("product not found")
	}
	m.preloadCategory(ctx, p)
	return p, nil
}

func (m *MockProductRepository) FindBySKU(_ context.Context, sku string) (*models.Product, error) {
	for _, p := range m.Products {
		if isActiveProduct(p) && p.SKU == sku {
			return p, nil
		}
	}
	return nil, errors.New("product not found")
}

func (m *MockProductRepository) Create(_ context.Context, product *models.Product) error {
	product.ID = m.NextID
	m.NextID++
	m.Products[product.ID] = product
	return nil
}

func (m *MockProductRepository) Update(_ context.Context, product *models.Product) error {
	p, ok := m.Products[product.ID]
	if !ok || !isActiveProduct(p) {
		return errors.New("product not found")
	}
	m.Products[product.ID] = product
	return nil
}

func (m *MockProductRepository) Delete(_ context.Context, id uint) error {
	p, ok := m.Products[id]
	if !ok || !isActiveProduct(p) {
		return errors.New("product not found")
	}
	p.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	return nil
}

func (m *MockProductRepository) CountByCategoryID(_ context.Context, categoryID uint) (int64, error) {
	var count int64
	for _, p := range m.Products {
		if isActiveProduct(p) && p.CategoryID == categoryID {
			count++
		}
	}
	return count, nil
}
