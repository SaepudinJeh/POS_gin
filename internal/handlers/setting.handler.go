package handlers

import (
	"POS/internal/dto"
	"POS/internal/response"
	"POS/internal/services"
	"POS/internal/validators"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	svc *services.SettingService
}

func NewSettingHandler(svc *services.SettingService) *SettingHandler {
	return &SettingHandler{svc: svc}
}

// GetAllSettings godoc
// @Summary      List semua setting (admin only)
// @Tags         settings
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response
// @Failure      401 {object} response.ErrorResponse
// @Failure      403 {object} response.ErrorResponse
// @Router       /settings [get]
func (h *SettingHandler) GetAll(c *gin.Context) {
	settings, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	result := make([]dto.SettingResponse, 0, len(settings))
	for _, s := range settings {
		result = append(result, dto.SettingResponse{
			Key:         s.Key,
			Value:       s.Value,
			Description: s.Description,
		})
	}
	response.Success(c, result)
}

// UpdateSetting godoc
// @Summary      Update satu setting (admin only)
// @Tags         settings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        key path string true "Key setting (contoh: tax_rate)"
// @Param        request body dto.UpdateSettingRequest true "Nilai baru"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      403 {object} response.ErrorResponse
// @Router       /settings/{key} [put]
func (h *SettingHandler) Update(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "Key tidak valid")
		return
	}

	var req dto.UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validators.FormatValidationError(err))
		return
	}

	if err := h.svc.Update(c.Request.Context(), key, req.Value); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "Setting berhasil diperbarui", nil)
}

// BulkUpdateSettings godoc
// @Summary      Update banyak setting sekaligus (admin only)
// @Tags         settings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.BulkUpdateSettingRequest true "Map key-value"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      403 {object} response.ErrorResponse
// @Router       /settings [put]
func (h *SettingHandler) BulkUpdate(c *gin.Context) {
	var req dto.BulkUpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validators.FormatValidationError(err))
		return
	}

	if err := h.svc.BulkUpdate(c.Request.Context(), req.Settings); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "Settings berhasil diperbarui", nil)
}
