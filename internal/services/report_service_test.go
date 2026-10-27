package services_test

import (
	"context"
	"testing"
	"time"

	"POS/internal/models"
	"POS/internal/repositories"
	"POS/internal/services"
)

func TestParseRange_Defaults(t *testing.T) {
	from, to, err := services.ParseRange("", "")
	if err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if to.Sub(from) != 24*time.Hour {
		t.Errorf("range default harus 1 hari, dapat: %v", to.Sub(from))
	}
}

func TestParseRange_Valid(t *testing.T) {
	from, to, err := services.ParseRange("2026-10-01", "2026-10-05")
	if err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if to.Sub(from) != 5*24*time.Hour {
		t.Errorf("range harus 5 hari, dapat: %v", to.Sub(from))
	}
}

func TestParseRange_InvalidFormat(t *testing.T) {
	_, _, err := services.ParseRange("01-10-2026", "")
	if err == nil {
		t.Fatal("harusnya error format")
	}
}

func TestParseRange_ToBeforeFrom(t *testing.T) {
	_, _, err := services.ParseRange("2026-10-05", "2026-10-01")
	if err == nil {
		t.Fatal("harusnya error to < from")
	}
}

func TestReportService_Summary(t *testing.T) {
	ctx := context.Background()
	repo := repositories.NewMockReportRepository()
	svc := services.NewReportService(repo)

	now := time.Now()
	repo.Transactions = []*models.Transaction{
		{
			Status:         models.TransactionStatusPaid,
			Total:          10000,
			DiscountAmount: 500,
			TaxAmount:      1000,
			CreatedAt:      now,
			Items: []*models.TransactionItem{
				{UnitPrice: 5000, UnitCost: 3000, Quantity: 2},
			},
		},
		{
			Status:    models.TransactionStatusPending,
			Total:     20000,
			CreatedAt: now,
		},
	}

	from, to, _ := services.ParseRange("", "")
	data, err := svc.Summary(ctx, from, to)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if data.TotalTransactions != 1 {
		t.Errorf("harusnya 1, dapat: %d", data.TotalTransactions)
	}
	if data.TotalRevenue != 10000 {
		t.Errorf("revenue harusnya 10000, dapat: %f", data.TotalRevenue)
	}
	if data.TotalProfit != 4000 {
		t.Errorf("profit harusnya 4000, dapat: %f", data.TotalProfit)
	}
}
