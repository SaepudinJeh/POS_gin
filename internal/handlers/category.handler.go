package handlers

import (
	"strconv"

	"POS/internal/dto"
	"POS/internal/repositories"
	"POS/internal/response"
	"POS/internal/services"
	"POS/internal/validators"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	svc *services.CategoryService
}

func NewCategoryHandler(svc *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// GetAllCategories godoc
// @Summary      List kategori
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Halaman (default 1)"
// @Param        limit query int false "Jumlah per halaman (default 10)"
// @Param        search query string false "Cari berdasarkan nama"
// @Success      200 {object} response.PaginationResponse
// @Failure      401 {object} response.ErrorResponse
// @Router       /categories [get]
func (h *CategoryHandler) GetAll(c *gin.Context) {
	filter := repositories.CategoryFilter{
		Search: c.Query("search"),
		Page:   parseIntDefault(c.Query("page"), 1),
		Limit:  parseIntDefault(c.Query("limit"), 10),
	}

	categories, total, err := h.svc.GetAll(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Paginated(c, categories, total, filter.Page, filter.Limit)
}

// GetCategory godoc
// @Summary      Detail kategori
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID kategori"
// @Success      200 {object} response.Response
// @Failure      404 {object} response.ErrorResponse
// @Router       /categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	category, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, category)
}

// CreateCategory godoc
// @Summary      Buat kategori baru
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateCategoryRequest true "Data kategori"
// @Success      201 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      403 {object} response.ErrorResponse
// @Router       /categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validators.FormatValidationError(err))
		return
	}

	category, err := h.svc.Create(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.CreatedWithMessage(c, "Kategori berhasil dibuat", category)
}

// UpdateCategory godoc
// @Summary      Update kategori
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID kategori"
// @Param        request body dto.UpdateCategoryRequest true "Data kategori"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      403 {object} response.ErrorResponse
// @Router       /categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validators.FormatValidationError(err))
		return
	}

	category, err := h.svc.Update(c.Request.Context(), uint(id), req.Name, req.Description)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "Kategori berhasil diperbarui", category)
}

// DeleteCategory godoc
// @Summary      Hapus kategori (soft delete)
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID kategori"
// @Success      200 {object} response.Response
// @Failure      401 {object} response.ErrorResponse
// @Failure      403 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "Kategori berhasil dihapus", nil)
}
