package handlers

import (
	"strconv"
	"time"

	"POS/internal/dto"
	"POS/internal/models"
	"POS/internal/repositories"
	"POS/internal/response"
	"POS/internal/services"
	"POS/internal/validators"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	svc *services.TransactionService
}

func NewTransactionHandler(svc *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

func (h *TransactionHandler) GetAll(c *gin.Context) {
	filter := repositories.TransactionFilter{
		Status: c.Query("status"),
		Search: c.Query("search"),
		Page:   parseIntDefault(c.Query("page"), 1),
		Limit:  parseIntDefault(c.Query("limit"), 10),
	}

	if cidStr := c.Query("cashier_id"); cidStr != "" {
		if id, err := strconv.ParseUint(cidStr, 10, 64); err == nil {
			cid := uint(id)
			filter.CashierID = &cid
		}
	}

	if from := c.Query("from"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			filter.From = &t
		}
	}
	if to := c.Query("to"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			endOfDay := t.Add(24 * time.Hour)
			filter.To = &endOfDay
		}
	}

	txns, total, err := h.svc.GetAll(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	responses := make([]dto.TransactionResponse, 0, len(txns))
	for _, t := range txns {
		responses = append(responses, toTransactionResponse(t))
	}

	response.Paginated(c, responses, total, filter.Page, filter.Limit)
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	txn, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, toTransactionResponse(txn))
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validators.FormatValidationError(err))
		return
	}

	userID := getUserID(c)

	items := make([]services.TransactionItemInput, 0, len(req.Items))
	for _, in := range req.Items {
		items = append(items, services.TransactionItemInput{
			ProductID: in.ProductID,
			Quantity:  in.Quantity,
		})
	}

	txn, err := h.svc.Create(c.Request.Context(), services.TransactionInput{
		Items:         items,
		DiscountType:  req.DiscountType,
		DiscountValue: req.DiscountValue,
		PaymentMethod: req.PaymentMethod,
		PaymentAmount: req.PaymentAmount,
		Note:          req.Note,
		CashierID:     userID,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.CreatedWithMessage(c, "Transaksi berhasil dibuat", toTransactionResponse(txn))
}

func (h *TransactionHandler) Pay(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	var req dto.PayTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validators.FormatValidationError(err))
		return
	}

	txn, err := h.svc.Pay(c.Request.Context(), uint(id), req.PaymentMethod, req.PaymentAmount)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "Pembayaran berhasil", toTransactionResponse(txn))
}

func (h *TransactionHandler) Void(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	var req dto.VoidTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validators.FormatValidationError(err))
		return
	}

	txn, err := h.svc.Void(c.Request.Context(), uint(id), req.Reason, getUserID(c))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "Transaksi berhasil di-void", toTransactionResponse(txn))
}

func getUserID(c *gin.Context) uint {
	if v, ok := c.Get("user_id"); ok {
		if f, ok := v.(float64); ok {
			return uint(f)
		}
	}
	return 0
}

func toTransactionResponse(t *models.Transaction) dto.TransactionResponse {
	resp := dto.TransactionResponse{
		ID:             t.ID,
		InvoiceNumber:  t.InvoiceNumber,
		CashierID:      t.CashierID,
		Status:         t.Status,
		Subtotal:       t.Subtotal,
		DiscountType:   t.DiscountType,
		DiscountValue:  t.DiscountValue,
		DiscountAmount: t.DiscountAmount,
		TaxRate:        t.TaxRate,
		TaxAmount:      t.TaxAmount,
		Total:          t.Total,
		PaymentMethod:  t.PaymentMethod,
		PaymentAmount:  t.PaymentAmount,
		ChangeAmount:   t.ChangeAmount,
		PaidAt:         t.PaidAt,
		VoidedAt:       t.VoidedAt,
		VoidedBy:       t.VoidedBy,
		VoidReason:     t.VoidReason,
		Note:           t.Note,
		CreatedAt:      t.CreatedAt,
	}
	for _, item := range t.Items {
		resp.Items = append(resp.Items, dto.TransactionItemResponse{
			ProductID:   item.ProductID,
			ProductSKU:  item.ProductSKU,
			ProductName: item.ProductName,
			UnitPrice:   item.UnitPrice,
			Quantity:    item.Quantity,
			Subtotal:    item.Subtotal,
		})
	}
	if resp.Items == nil {
		resp.Items = []dto.TransactionItemResponse{}
	}
	return resp
}
