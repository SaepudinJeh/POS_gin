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

type ProductHandler struct {
	svc *services.ProductService
}

func NewProductHandler(svc *services.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	filter := repositories.ProductFilter{
		Search: c.Query("search"),
		Page:   parseIntDefault(c.Query("page"), 1),
		Limit:  parseIntDefault(c.Query("limit"), 10),
	}

	if catIDStr := c.Query("category_id"); catIDStr != "" {
		if id, err := strconv.ParseUint(catIDStr, 10, 64); err == nil {
			catID := uint(id)
			filter.CategoryID = &catID
		}
	}

	if activeStr := c.Query("is_active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			filter.IsActive = &active
		}
	}

	products, total, err := h.svc.GetAll(filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	responses := make([]dto.ProductResponse, 0, len(products))
	for _, p := range products {
		responses = append(responses, toProductResponse(p))
	}

	response.Paginated(c, responses, total, filter.Page, filter.Limit)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	product, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, toProductResponse(product))
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validators.FormatValidationError(err))
		return
	}

	product, err := h.svc.Create(services.ProductInput{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Cost:        req.Cost,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	full, _ := h.svc.GetByID(product.ID)
	response.CreatedWithMessage(c, "Produk berhasil dibuat", toProductResponse(full))
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validators.FormatValidationError(err))
		return
	}

	product, err := h.svc.Update(uint(id), services.ProductInput{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Cost:        req.Cost,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		IsActive:    req.IsActive,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	full, _ := h.svc.GetByID(product.ID)
	response.SuccessWithMessage(c, "Produk berhasil diperbarui", toProductResponse(full))
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "Produk berhasil dihapus", nil)
}

func toProductResponse(p *models.Product) dto.ProductResponse {
	if p == nil {
		return dto.ProductResponse{}
	}
	resp := dto.ProductResponse{
		ID:          p.ID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Cost:        p.Cost,
		Stock:       p.Stock,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	if p.Category != nil {
		resp.Category = &dto.CategoryResponse{
			ID:          p.Category.ID,
			Name:        p.Category.Name,
			Description: p.Category.Description,
		}
	}
	return resp
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
