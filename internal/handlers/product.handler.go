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

// GetAllProducts godoc
// @Summary      List produk
// @Description  Filter & pagination. Contoh: ?page=1&limit=10&search=indomie&category_id=1&is_active=true
// @Tags         products
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Halaman (default 1)"
// @Param        limit query int false "Jumlah per halaman (default 10)"
// @Param        search query string false "Cari berdasarkan nama atau SKU"
// @Param        category_id query int false "Filter berdasarkan ID kategori"
// @Param        is_active query bool false "Filter status aktif"
// @Success      200 {object} response.PaginationResponse
// @Failure      401 {object} response.ErrorResponse
// @Router       /products [get]
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

	products, total, err := h.svc.GetAll(c.Request.Context(), filter)
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

// GetProduct godoc
// @Summary      Detail produk
// @Tags         products
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID produk"
// @Success      200 {object} response.Response
// @Failure      401 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	product, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, toProductResponse(product))
}

// CreateProduct godoc
// @Summary      Buat produk baru
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateProductRequest true "Data produk"
// @Success      201 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      403 {object} response.ErrorResponse
// @Router       /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validators.FormatValidationError(err))
		return
	}

	product, err := h.svc.Create(c.Request.Context(), services.ProductInput{
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

	full, _ := h.svc.GetByID(c.Request.Context(), product.ID)
	response.CreatedWithMessage(c, "Produk berhasil dibuat", toProductResponse(full))
}

// UpdateProduct godoc
// @Summary      Update produk
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID produk"
// @Param        request body dto.UpdateProductRequest true "Data produk"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      403 {object} response.ErrorResponse
// @Router       /products/{id} [put]
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

	product, err := h.svc.Update(c.Request.Context(), uint(id), services.ProductInput{
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

	full, _ := h.svc.GetByID(c.Request.Context(), product.ID)
	response.SuccessWithMessage(c, "Produk berhasil diperbarui", toProductResponse(full))
}

// DeleteProduct godoc
// @Summary      Hapus produk (soft delete)
// @Tags         products
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID produk"
// @Success      200 {object} response.Response
// @Failure      401 {object} response.ErrorResponse
// @Failure      403 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
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
