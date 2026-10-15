package handlers

import (
	"strconv"

	"POS/internal/dto"
	"POS/internal/models"
	"POS/internal/repositories"
	"POS/internal/response"
	"POS/internal/services"
	"POS/internal/validators"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	svc *services.StockService
}

func NewStockHandler(svc *services.StockService) *StockHandler {
	return &StockHandler{svc: svc}
}

func (h *StockHandler) GetAll(c *gin.Context) {
	filter := repositories.StockFilter{
		Type:  c.Query("type"),
		Page:  parseIntDefault(c.Query("page"), 1),
		Limit: parseIntDefault(c.Query("limit"), 10),
	}

	if pidStr := c.Query("product_id"); pidStr != "" {
		if id, err := strconv.ParseUint(pidStr, 10, 64); err == nil {
			pid := uint(id)
			filter.ProductID = &pid
		}
	}

	mutations, total, err := h.svc.GetAll(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	responses := make([]dto.StockMutationResponse, 0, len(mutations))
	for _, m := range mutations {
		responses = append(responses, toStockMutationResponse(m))
	}

	response.Paginated(c, responses, total, filter.Page, filter.Limit)
}

func (h *StockHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	mut, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, toStockMutationResponse(mut))
}

func (h *StockHandler) GetByProductID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	mutations, err := h.svc.GetByProductID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	responses := make([]dto.StockMutationResponse, 0, len(mutations))
	for _, m := range mutations {
		responses = append(responses, toStockMutationResponse(m))
	}

	response.Success(c, responses)
}

func (h *StockHandler) Create(c *gin.Context) {
	var req dto.CreateStockMutationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validators.FormatValidationError(err))
		return
	}

	userIDVal, _ := c.Get("user_id")
	var userID uint
	if f, ok := userIDVal.(float64); ok {
		userID = uint(f)
	}

	mut, err := h.svc.CreateMutation(c.Request.Context(), services.StockInput{
		ProductID: req.ProductID,
		Type:      req.Type,
		Quantity:  req.Quantity,
		Note:      req.Note,
		UserID:    userID,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.CreatedWithMessage(c, "Mutasi stok berhasil dicatat", toStockMutationResponse(mut))
}

func toStockMutationResponse(m *models.StockMutation) dto.StockMutationResponse {
	resp := dto.StockMutationResponse{
		ID:          m.ID,
		ProductID:   m.ProductID,
		Type:        m.Type,
		Quantity:    m.Quantity,
		StockBefore: m.StockBefore,
		StockAfter:  m.StockAfter,
		Note:        m.Note,
		UserID:      m.UserID,
		CreatedAt:   m.CreatedAt,
	}
	if m.Product != nil {
		p := toProductResponse(m.Product)
		resp.Product = &p
	}
	return resp
}
