package services

import (
	"context"
	"errors"
	"time"

	"POS/internal/repositories"
)

type ReportService struct {
	repo repositories.ReportRepository
}

func NewReportService(repo repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

// ParseRange mengambil from & to dari string, default hari ini kalau kosong.
func ParseRange(fromStr, toStr string) (time.Time, time.Time, error) {
	now := time.Now()
	loc := now.Location()

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	defaultFrom := today
	defaultTo := today.Add(24 * time.Hour)

	var from, to time.Time
	var err error

	if fromStr == "" {
		from = defaultFrom
	} else {
		from, err = time.ParseInLocation("2006-01-02", fromStr, loc)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("format 'from' harus YYYY-MM-DD")
		}
	}

	if toStr == "" {
		to = defaultTo
	} else {
		t, err := time.ParseInLocation("2006-01-02", toStr, loc)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("format 'to' harus YYYY-MM-DD")
		}
		to = t.Add(24 * time.Hour)
	}

	if !to.After(from) {
		return time.Time{}, time.Time{}, errors.New("'to' harus lebih besar dari 'from'")
	}

	return from, to, nil
}

func (s *ReportService) Summary(ctx context.Context, from, to time.Time) (*repositories.SummaryData, error) {
	return s.repo.Summary(ctx, from, to)
}

func (s *ReportService) Daily(ctx context.Context, from, to time.Time) ([]*repositories.DailyData, error) {
	return s.repo.Daily(ctx, from, to)
}

func (s *ReportService) ByCashier(ctx context.Context, from, to time.Time) ([]*repositories.CashierData, error) {
	return s.repo.ByCashier(ctx, from, to)
}

func (s *ReportService) ByProduct(ctx context.Context, from, to time.Time, limit int) ([]*repositories.ProductData, error) {
	return s.repo.ByProduct(ctx, from, to, limit)
}

func (s *ReportService) ByCategory(ctx context.Context, from, to time.Time) ([]*repositories.CategoryData, error) {
	return s.repo.ByCategory(ctx, from, to)
}
