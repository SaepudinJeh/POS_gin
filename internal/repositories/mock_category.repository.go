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

type MockCategoryRepository struct {
	Categories map[uint]*models.Category
	NextID     uint
}

func NewMockCategoryRepository() *MockCategoryRepository {
	return &MockCategoryRepository{
		Categories: make(map[uint]*models.Category),
		NextID:     1,
	}
}

func isActiveCategory(category *models.Category) bool {
	return !category.DeletedAt.Valid
}

func (m *MockCategoryRepository) FindAll(_ context.Context, filter CategoryFilter) ([]*models.Category, int64, error) {
	var result []*models.Category

	for _, cat := range m.Categories {
		if !isActiveCategory(cat) {
			continue
		}
		if filter.Search != "" {
			if !strings.Contains(strings.ToLower(cat.Name), strings.ToLower(filter.Search)) {
				continue
			}
		}
		result = append(result, cat)
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
		return []*models.Category{}, total, nil
	}
	end := offset + filter.Limit
	if end > len(result) {
		end = len(result)
	}

	return result[offset:end], total, nil
}

func (m *MockCategoryRepository) FindByID(_ context.Context, id uint) (*models.Category, error) {
	c, ok := m.Categories[id]
	if !ok || !isActiveCategory(c) {
		return nil, errors.New("category not found")
	}
	return c, nil
}

func (m *MockCategoryRepository) FindByName(_ context.Context, name string) (*models.Category, error) {
	for _, c := range m.Categories {
		if isActiveCategory(c) && c.Name == name {
			return c, nil
		}
	}
	return nil, errors.New("category not found")
}

func (m *MockCategoryRepository) Create(_ context.Context, category *models.Category) error {
	category.ID = m.NextID
	m.NextID++
	m.Categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepository) Update(_ context.Context, category *models.Category) error {
	c, ok := m.Categories[category.ID]
	if !ok || !isActiveCategory(c) {
		return errors.New("category not found")
	}
	m.Categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepository) Delete(_ context.Context, id uint) error {
	c, ok := m.Categories[id]
	if !ok || !isActiveCategory(c) {
		return errors.New("category not found")
	}
	c.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	return nil
}
