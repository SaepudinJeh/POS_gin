package dto

type PeriodInfo struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type ReportSummary struct {
	Period             PeriodInfo `json:"period"`
	TotalTransactions  int64      `json:"total_transactions"`
	TotalRevenue       float64    `json:"total_revenue"`
	TotalProfit        float64    `json:"total_profit"`
	TotalDiscount      float64    `json:"total_discount"`
	TotalTax           float64    `json:"total_tax"`
	AverageTransaction float64    `json:"average_transaction"`
}

type DailyReportItem struct {
	Date    string  `json:"date"`
	Count   int64   `json:"count"`
	Revenue float64 `json:"revenue"`
}

type CashierReportItem struct {
	CashierID   uint    `json:"cashier_id"`
	CashierName string  `json:"cashier_name"`
	Count       int64   `json:"count"`
	Revenue     float64 `json:"revenue"`
}

type ProductReportItem struct {
	ProductID   uint    `json:"product_id"`
	ProductSKU  string  `json:"product_sku"`
	ProductName string  `json:"product_name"`
	Qty         int64   `json:"qty"`
	Revenue     float64 `json:"revenue"`
	Profit      float64 `json:"profit"`
}

type CategoryReportItem struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Qty          int64   `json:"qty"`
	Revenue      float64 `json:"revenue"`
}
