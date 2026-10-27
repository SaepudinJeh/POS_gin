package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type SummaryData struct {
	TotalTransactions int64
	TotalRevenue      float64
	TotalDiscount     float64
	TotalTax          float64
	TotalProfit       float64
}

type DailyData struct {
	Date    string
	Count   int64
	Revenue float64
}

type CashierData struct {
	CashierID   uint
	CashierName string
	Count       int64
	Revenue     float64
}

type ProductData struct {
	ProductID   uint
	ProductSKU  string
	ProductName string
	Qty         int64
	Revenue     float64
	Profit      float64
}

type CategoryData struct {
	CategoryID   uint
	CategoryName string
	Qty          int64
	Revenue      float64
}

type ReportRepository interface {
	Summary(ctx context.Context, from, to time.Time) (*SummaryData, error)
	Daily(ctx context.Context, from, to time.Time) ([]*DailyData, error)
	ByCashier(ctx context.Context, from, to time.Time) ([]*CashierData, error)
	ByProduct(ctx context.Context, from, to time.Time, limit int) ([]*ProductData, error)
	ByCategory(ctx context.Context, from, to time.Time) ([]*CategoryData, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) Summary(ctx context.Context, from, to time.Time) (*SummaryData, error) {
	// Total dari transaksi paid
	var txRow struct {
		TotalTransactions int64
		TotalRevenue      float64
		TotalDiscount     float64
		TotalTax          float64
	}

	if err := r.db.WithContext(ctx).
		Raw(`SELECT
			COUNT(*) AS total_transactions,
			COALESCE(SUM(total), 0) AS total_revenue,
			COALESCE(SUM(discount_amount), 0) AS total_discount,
			COALESCE(SUM(tax_amount), 0) AS total_tax
			FROM transactions
			WHERE status = 'paid' AND deleted_at IS NULL
			  AND created_at >= ? AND created_at < ?`,
			from, to).
		Scan(&txRow).Error; err != nil {
		return nil, err
	}

	// Total profit dari items
	var profitRow struct {
		TotalProfit float64
	}
	if err := r.db.WithContext(ctx).
		Raw(`SELECT COALESCE(SUM((ti.unit_price - ti.unit_cost) * ti.quantity), 0) AS total_profit
			FROM transaction_items ti
			JOIN transactions t ON t.id = ti.transaction_id
			WHERE t.status = 'paid' AND t.deleted_at IS NULL
			  AND t.created_at >= ? AND t.created_at < ?`,
			from, to).
		Scan(&profitRow).Error; err != nil {
		return nil, err
	}

	return &SummaryData{
		TotalTransactions: txRow.TotalTransactions,
		TotalRevenue:      txRow.TotalRevenue,
		TotalDiscount:     txRow.TotalDiscount,
		TotalTax:          txRow.TotalTax,
		TotalProfit:       profitRow.TotalProfit,
	}, nil
}

func (r *reportRepository) Daily(ctx context.Context, from, to time.Time) ([]*DailyData, error) {
	var rows []*DailyData
	if err := r.db.WithContext(ctx).
		Raw(`SELECT
			TO_CHAR(created_at, 'YYYY-MM-DD') AS date,
			COUNT(*) AS count,
			COALESCE(SUM(total), 0) AS revenue
			FROM transactions
			WHERE status = 'paid' AND deleted_at IS NULL
			  AND created_at >= ? AND created_at < ?
			GROUP BY TO_CHAR(created_at, 'YYYY-MM-DD')
			ORDER BY date ASC`,
			from, to).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *reportRepository) ByCashier(ctx context.Context, from, to time.Time) ([]*CashierData, error) {
	var rows []*CashierData
	if err := r.db.WithContext(ctx).
		Raw(`SELECT
			t.cashier_id,
			COALESCE(u.name, 'Unknown') AS cashier_name,
			COUNT(*) AS count,
			COALESCE(SUM(t.total), 0) AS revenue
			FROM transactions t
			LEFT JOIN users u ON u.id = t.cashier_id
			WHERE t.status = 'paid' AND t.deleted_at IS NULL
			  AND t.created_at >= ? AND t.created_at < ?
			GROUP BY t.cashier_id, u.name
			ORDER BY revenue DESC`,
			from, to).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *reportRepository) ByProduct(ctx context.Context, from, to time.Time, limit int) ([]*ProductData, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var rows []*ProductData
	if err := r.db.WithContext(ctx).
		Raw(`SELECT
			ti.product_id,
			ti.product_sku,
			ti.product_name,
			SUM(ti.quantity) AS qty,
			COALESCE(SUM(ti.subtotal), 0) AS revenue,
			COALESCE(SUM((ti.unit_price - ti.unit_cost) * ti.quantity), 0) AS profit
			FROM transaction_items ti
			JOIN transactions t ON t.id = ti.transaction_id
			WHERE t.status = 'paid' AND t.deleted_at IS NULL
			  AND t.created_at >= ? AND t.created_at < ?
			GROUP BY ti.product_id, ti.product_sku, ti.product_name
			ORDER BY revenue DESC
			LIMIT ?`,
			from, to, limit).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *reportRepository) ByCategory(ctx context.Context, from, to time.Time) ([]*CategoryData, error) {
	var rows []*CategoryData
	if err := r.db.WithContext(ctx).
		Raw(`SELECT
			COALESCE(p.category_id, 0) AS category_id,
			COALESCE(c.name, 'Tanpa Kategori') AS category_name,
			SUM(ti.quantity) AS qty,
			COALESCE(SUM(ti.subtotal), 0) AS revenue
			FROM transaction_items ti
			JOIN transactions t ON t.id = ti.transaction_id
			LEFT JOIN products p ON p.id = ti.product_id
			LEFT JOIN categories c ON c.id = p.category_id
			WHERE t.status = 'paid' AND t.deleted_at IS NULL
			  AND t.created_at >= ? AND t.created_at < ?
			GROUP BY p.category_id, c.name
			ORDER BY revenue DESC`,
			from, to).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
