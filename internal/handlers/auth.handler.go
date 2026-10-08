package handlers

import (
	"POS/internal/dto"
	"POS/internal/response"
	"POS/internal/services"
	"POS/internal/validators"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validators.FormatValidationError(err))
		return
	}

	user, err := h.svc.Register(req.Name, req.Email, req.Password)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.CreatedWithMessage(c, "Registrasi berhasil", user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validators.FormatValidationError(err))
		return
	}

	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Login berhasil", gin.H{"token": token})
}
