package handlers

import (
	"strconv"

	"POS/internal/dto"
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

func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.svc.GetAll()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, categories)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	category, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, category)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validators.FormatValidationError(err))
		return
	}

	category, err := h.svc.Create(req.Name, req.Description)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.CreatedWithMessage(c, "Kategori berhasil dibuat", category)
}

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

	category, err := h.svc.Update(uint(id), req.Name, req.Description)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "Kategori berhasil diperbarui", category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID tidak valid")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "Kategori berhasil dihapus", nil)
}
