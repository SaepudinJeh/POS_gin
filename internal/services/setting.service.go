package services

import (
	"context"
	"strconv"

	"POS/internal/models"
	"POS/internal/repositories"
)

const (
	SettingKeyTaxRate       = "tax_rate"
	SettingKeyStoreName     = "store_name"
	SettingKeyStoreAddress  = "store_address"
	SettingKeyStorePhone    = "store_phone"
	SettingKeyInvoicePrefix = "invoice_prefix"
)

type SettingService struct {
	repo repositories.SettingRepository
}

func NewSettingService(repo repositories.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

func (s *SettingService) GetAll(ctx context.Context) ([]*models.Setting, error) {
	return s.repo.FindAll(ctx)
}

func (s *SettingService) GetTaxRate(ctx context.Context) float64 {
	s2, err := s.repo.FindByKey(ctx, SettingKeyTaxRate)
	if err != nil {
		return 0
	}
	rate, err := strconv.ParseFloat(s2.Value, 64)
	if err != nil {
		return 0
	}
	return rate
}

func (s *SettingService) GetString(ctx context.Context, key, def string) string {
	s2, err := s.repo.FindByKey(ctx, key)
	if err != nil {
		return def
	}
	return s2.Value
}

func (s *SettingService) Update(ctx context.Context, key, value string) error {
	return s.repo.Upsert(ctx, key, value)
}

func (s *SettingService) BulkUpdate(ctx context.Context, settings map[string]string) error {
	return s.repo.BulkUpsert(ctx, settings)
}
