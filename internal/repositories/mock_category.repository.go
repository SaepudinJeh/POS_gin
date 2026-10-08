package repositories

import (
	"errors"
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

func (m *MockCategoryRepository) FindAll() ([]*models.Category, error) { // <-- pointer
	var result []*models.Category // <-- pointer
	for _, category := range m.Categories {
		if isActiveCategory(category) {
			result = append(result, category)
		}
	}
	return result, nil
}

func (m *MockCategoryRepository) FindByID(id uint) (*models.Category, error) {
	c, ok := m.Categories[id]
	if !ok || !isActiveCategory(c) {
		return nil, errors.New("category not found")
	}
	return c, nil
}

func (m *MockCategoryRepository) FindByName(name string) (*models.Category, error) {
	for _, c := range m.Categories {
		if isActiveCategory(c) && c.Name == name {
			return c, nil
		}
	}
	return nil, errors.New("category not found")
}

func (m *MockCategoryRepository) Create(category *models.Category) error {
	category.ID = m.NextID
	m.NextID++
	m.Categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepository) Update(category *models.Category) error {
	c, ok := m.Categories[category.ID]
	if !ok || !isActiveCategory(c) {
		return errors.New("category not found")
	}
	m.Categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepository) Delete(id uint) error {
	c, ok := m.Categories[id]
	if !ok || !isActiveCategory(c) {
		return errors.New("category not found")
	}
	c.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	return nil
}
