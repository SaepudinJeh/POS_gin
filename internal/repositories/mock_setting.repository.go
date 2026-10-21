package repositories

import (
	"context"
	"errors"

	"POS/internal/models"
)

type MockSettingRepository struct {
	Settings map[string]*models.Setting
}

func NewMockSettingRepository() *MockSettingRepository {
	return &MockSettingRepository{Settings: make(map[string]*models.Setting)}
}

func (m *MockSettingRepository) FindAll(_ context.Context) ([]*models.Setting, error) {
	var result []*models.Setting
	for _, s := range m.Settings {
		result = append(result, s)
	}
	return result, nil
}

func (m *MockSettingRepository) FindByKey(_ context.Context, key string) (*models.Setting, error) {
	s, ok := m.Settings[key]
	if !ok {
		return nil, errors.New("setting not found")
	}
	return s, nil
}

func (m *MockSettingRepository) Upsert(_ context.Context, key, value string) error {
	if s, ok := m.Settings[key]; ok {
		s.Value = value
		return nil
	}
	m.Settings[key] = &models.Setting{Key: key, Value: value}
	return nil
}

func (m *MockSettingRepository) BulkUpsert(ctx context.Context, settings map[string]string) error {
	for k, v := range settings {
		if err := m.Upsert(ctx, k, v); err != nil {
			return err
		}
	}
	return nil
}
