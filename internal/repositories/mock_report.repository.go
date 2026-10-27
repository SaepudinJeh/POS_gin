package repositories

import (
	"context"
	"sort"
	"time"

	"POS/internal/models"
)

type MockReportRepository struct {
	Transactions []*models.Transaction
	Users        map[uint]*models.User
	Products     map[uint]*models.Product
	Categories   map[uint]*models.Category
}

func NewMockReportRepository() *MockReportRepository {
	return &MockReportRepository{
		Users:      make(map[uint]*models.User),
		Products:   make(map[uint]*models.Product),
		Categories: make(map[uint]*models.Category),
	}
}

func inRange(t time.Time, from, to time.Time) bool {
	return !t.Before(from) && t.Before(to)
}

func isPaidTxn(t *models.Transaction) bool {
	return t.Status == models.TransactionStatusPaid && !t.DeletedAt.Valid
}

func (m *MockReportRepository) Summary(_ context.Context, from, to time.Time) (*SummaryData, error) {
	row := &SummaryData{}
	for _, t := range m.Transactions {
		if !isPaidTxn(t) || !inRange(t.CreatedAt, from, to) {
			continue
		}
		row.TotalTransactions++
		row.TotalRevenue += t.Total
		row.TotalDiscount += t.DiscountAmount
		row.TotalTax += t.TaxAmount
		for _, item := range t.Items {
			row.TotalProfit += (item.UnitPrice - item.UnitCost) * float64(item.Quantity)
		}
	}
	return row, nil
}

func (m *MockReportRepository) Daily(_ context.Context, from, to time.Time) ([]*DailyData, error) {
	group := make(map[string]*DailyData)
	for _, t := range m.Transactions {
		if !isPaidTxn(t) || !inRange(t.CreatedAt, from, to) {
			continue
		}
		date := t.CreatedAt.Format("2006-01-02")
		if _, ok := group[date]; !ok {
			group[date] = &DailyData{Date: date}
		}
		group[date].Count++
		group[date].Revenue += t.Total
	}
	result := make([]*DailyData, 0, len(group))
	for _, v := range group {
		result = append(result, v)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Date < result[j].Date })
	return result, nil
}

func (m *MockReportRepository) ByCashier(_ context.Context, from, to time.Time) ([]*CashierData, error) {
	group := make(map[uint]*CashierData)
	for _, t := range m.Transactions {
		if !isPaidTxn(t) || !inRange(t.CreatedAt, from, to) {
			continue
		}
		if _, ok := group[t.CashierID]; !ok {
			name := "Unknown"
			if u, ok := m.Users[t.CashierID]; ok {
				name = u.Name
			}
			group[t.CashierID] = &CashierData{CashierID: t.CashierID, CashierName: name}
		}
		group[t.CashierID].Count++
		group[t.CashierID].Revenue += t.Total
	}
	result := make([]*CashierData, 0, len(group))
	for _, v := range group {
		result = append(result, v)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Revenue > result[j].Revenue })
	return result, nil
}

func (m *MockReportRepository) ByProduct(_ context.Context, from, to time.Time, limit int) ([]*ProductData, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}
	group := make(map[uint]*ProductData)
	for _, t := range m.Transactions {
		if !isPaidTxn(t) || !inRange(t.CreatedAt, from, to) {
			continue
		}
		for _, item := range t.Items {
			if _, ok := group[item.ProductID]; !ok {
				group[item.ProductID] = &ProductData{
					ProductID:   item.ProductID,
					ProductSKU:  item.ProductSKU,
					ProductName: item.ProductName,
				}
			}
			group[item.ProductID].Qty += int64(item.Quantity)
			group[item.ProductID].Revenue += item.Subtotal
			group[item.ProductID].Profit += (item.UnitPrice - item.UnitCost) * float64(item.Quantity)
		}
	}
	result := make([]*ProductData, 0, len(group))
	for _, v := range group {
		result = append(result, v)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Revenue > result[j].Revenue })
	if len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

func (m *MockReportRepository) ByCategory(_ context.Context, from, to time.Time) ([]*CategoryData, error) {
	group := make(map[uint]*CategoryData)
	for _, t := range m.Transactions {
		if !isPaidTxn(t) || !inRange(t.CreatedAt, from, to) {
			continue
		}
		for _, item := range t.Items {
			var catID uint
			var catName = "Tanpa Kategori"
			if p, ok := m.Products[item.ProductID]; ok {
				catID = p.CategoryID
				if c, ok := m.Categories[catID]; ok {
					catName = c.Name
				}
			}
			if _, ok := group[catID]; !ok {
				group[catID] = &CategoryData{CategoryID: catID, CategoryName: catName}
			}
			group[catID].Qty += int64(item.Quantity)
			group[catID].Revenue += item.Subtotal
		}
	}
	result := make([]*CategoryData, 0, len(group))
	for _, v := range group {
		result = append(result, v)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Revenue > result[j].Revenue })
	return result, nil
}
