package handlers

import (
	"POS/internal/dto"
	"POS/internal/response"
	"POS/internal/services"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	svc *services.ReportService
}

func NewReportHandler(svc *services.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

// Summary godoc
// @Summary      Ringkasan penjualan
// @Description  Total transaksi, revenue, profit, diskon, pajak. Hanya transaksi paid.
// @Tags         reports
// @Produce      json
// @Security     BearerAuth
// @Param        from query string false "Dari tanggal (YYYY-MM-DD), default hari ini"
// @Param        to query string false "Sampai tanggal (YYYY-MM-DD), default hari ini"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Router       /reports/summary [get]
func (h *ReportHandler) Summary(c *gin.Context) {
	from, to, err := services.ParseRange(c.Query("from"), c.Query("to"))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	data, err := h.svc.Summary(c.Request.Context(), from, to)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	var avg float64
	if data.TotalTransactions > 0 {
		avg = data.TotalRevenue / float64(data.TotalTransactions)
	}

	response.Success(c, dto.ReportSummary{
		Period:             periodInfo(from, to),
		TotalTransactions:  data.TotalTransactions,
		TotalRevenue:       data.TotalRevenue,
		TotalProfit:        data.TotalProfit,
		TotalDiscount:      data.TotalDiscount,
		TotalTax:           data.TotalTax,
		AverageTransaction: avg,
	})
}

// Daily godoc
// @Summary      Laporan harian
// @Description  Time series per hari untuk chart.
// @Tags         reports
// @Produce      json
// @Security     BearerAuth
// @Param        from query string false "Dari tanggal (YYYY-MM-DD)"
// @Param        to query string false "Sampai tanggal (YYYY-MM-DD)"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Router       /reports/daily [get]
func (h *ReportHandler) Daily(c *gin.Context) {
	from, to, err := services.ParseRange(c.Query("from"), c.Query("to"))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rows, err := h.svc.Daily(c.Request.Context(), from, to)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	result := make([]dto.DailyReportItem, 0, len(rows))
	for _, r := range rows {
		result = append(result, dto.DailyReportItem{
			Date:    r.Date,
			Count:   r.Count,
			Revenue: r.Revenue,
		})
	}
	response.Success(c, gin.H{
		"period": periodInfo(from, to),
		"data":   result,
	})
}

// ByCashier godoc
// @Summary      Laporan per kasir
// @Tags         reports
// @Produce      json
// @Security     BearerAuth
// @Param        from query string false "Dari tanggal (YYYY-MM-DD)"
// @Param        to query string false "Sampai tanggal (YYYY-MM-DD)"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Router       /reports/by-cashier [get]
func (h *ReportHandler) ByCashier(c *gin.Context) {
	from, to, err := services.ParseRange(c.Query("from"), c.Query("to"))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rows, err := h.svc.ByCashier(c.Request.Context(), from, to)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	result := make([]dto.CashierReportItem, 0, len(rows))
	for _, r := range rows {
		result = append(result, dto.CashierReportItem{
			CashierID:   r.CashierID,
			CashierName: r.CashierName,
			Count:       r.Count,
			Revenue:     r.Revenue,
		})
	}
	response.Success(c, gin.H{
		"period": periodInfo(from, to),
		"data":   result,
	})
}

// ByProduct godoc
// @Summary      Laporan produk terlaris
// @Tags         reports
// @Produce      json
// @Security     BearerAuth
// @Param        from query string false "Dari tanggal (YYYY-MM-DD)"
// @Param        to query string false "Sampai tanggal (YYYY-MM-DD)"
// @Param        limit query int false "Jumlah produk (default 10)"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Router       /reports/by-product [get]
func (h *ReportHandler) ByProduct(c *gin.Context) {
	from, to, err := services.ParseRange(c.Query("from"), c.Query("to"))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	limit := parseIntDefault(c.Query("limit"), 10)

	rows, err := h.svc.ByProduct(c.Request.Context(), from, to, limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	result := make([]dto.ProductReportItem, 0, len(rows))
	for _, r := range rows {
		result = append(result, dto.ProductReportItem{
			ProductID:   r.ProductID,
			ProductSKU:  r.ProductSKU,
			ProductName: r.ProductName,
			Qty:         r.Qty,
			Revenue:     r.Revenue,
			Profit:      r.Profit,
		})
	}
	response.Success(c, gin.H{
		"period": periodInfo(from, to),
		"data":   result,
	})
}

// ByCategory godoc
// @Summary      Laporan per kategori
// @Tags         reports
// @Produce      json
// @Security     BearerAuth
// @Param        from query string false "Dari tanggal (YYYY-MM-DD)"
// @Param        to query string false "Sampai tanggal (YYYY-MM-DD)"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Router       /reports/by-category [get]
func (h *ReportHandler) ByCategory(c *gin.Context) {
	from, to, err := services.ParseRange(c.Query("from"), c.Query("to"))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rows, err := h.svc.ByCategory(c.Request.Context(), from, to)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	result := make([]dto.CategoryReportItem, 0, len(rows))
	for _, r := range rows {
		result = append(result, dto.CategoryReportItem{
			CategoryID:   r.CategoryID,
			CategoryName: r.CategoryName,
			Qty:          r.Qty,
			Revenue:      r.Revenue,
		})
	}
	response.Success(c, gin.H{
		"period": periodInfo(from, to),
		"data":   result,
	})
}

func periodInfo(from, to interface{ Format(string) string }) dto.PeriodInfo {
	// to = batas eksklusif, kurangi 1 hari untuk tampilan
	return dto.PeriodInfo{
		From: from.Format("2006-01-02"),
		To:   to.Format("2006-01-02"),
	}
}
