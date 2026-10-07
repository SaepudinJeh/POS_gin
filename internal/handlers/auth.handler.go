package handlers

import (
	"POS/internal/dto"
	"POS/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	user, err := h.svc.Register(req.Name, req.Email, req.Password)

	if err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"success": true, "data": user})
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	token, err := h.svc.Login(req.Email, req.Password)

	if err != nil {
		ctx.JSON(401, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"success": true, "message": "Login successful", "token": token})
}
