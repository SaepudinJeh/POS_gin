package repositories

import (
	"errors"
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

func (m *MockProductRepository) preloadCategory(p *models.Product) {
	if m.CategoryRepo == nil {
		return
	}
	if cat, err := m.CategoryRepo.FindByID(p.CategoryID); err == nil {
		p.Category = cat
	}
}

func (m *MockProductRepository) FindAll(filter ProductFilter) ([]*models.Product, int64, error) {
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
		m.preloadCategory(p)
		result = append(result, p)
	}

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

func (m *MockProductRepository) FindByID(id uint) (*models.Product, error) {
	p, ok := m.Products[id]
	if !ok || !isActiveProduct(p) {
		return nil, errors.New("product not found")
	}
	m.preloadCategory(p)
	return p, nil
}

func (m *MockProductRepository) FindBySKU(sku string) (*models.Product, error) {
	for _, p := range m.Products {
		if isActiveProduct(p) && p.SKU == sku {
			return p, nil
		}
	}
	return nil, errors.New("product not found")
}

func (m *MockProductRepository) Create(product *models.Product) error {
	product.ID = m.NextID
	m.NextID++
	m.Products[product.ID] = product
	return nil
}

func (m *MockProductRepository) Update(product *models.Product) error {
	p, ok := m.Products[product.ID]
	if !ok || !isActiveProduct(p) {
		return errors.New("product not found")
	}
	m.Products[product.ID] = product
	return nil
}

func (m *MockProductRepository) Delete(id uint) error {
	p, ok := m.Products[id]
	if !ok || !isActiveProduct(p) {
		return errors.New("product not found")
	}
	p.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	return nil
}

func (m *MockProductRepository) CountByCategoryID(categoryID uint) (int64, error) {
	var count int64
	for _, p := range m.Products {
		if isActiveProduct(p) && p.CategoryID == categoryID {
			count++
		}
	}
	return count, nil
}
