package dto

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100" example:"Makanan"`
	Description string `json:"description" binding:"max=255" example:"Kategori makanan"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100" example:"Minuman"`
	Description string `json:"description" binding:"max=255" example:"Kategori minuman"`
}

type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
