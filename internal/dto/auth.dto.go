package dto

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3" example:"Budi Santoso"`
	Email    string `json:"email" binding:"required,email" example:"budi@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"rahasia123"`
}

type RegisterResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"budi@example.com"`
	Password string `json:"password" binding:"required" example:"rahasia123"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
