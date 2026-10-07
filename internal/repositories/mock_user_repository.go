package repositories

import (
	"POS/internal/models"

	"gorm.io/gorm"
)

type MockUserRepository struct {
	Users map[string]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	if user, exists := m.Users[email]; exists {
		return user, nil
	}

	return nil, gorm.ErrRecordNotFound
}

func (m *MockUserRepository) Create(user *models.User) error {
	m.Users[user.Email] = user
	return nil
}
