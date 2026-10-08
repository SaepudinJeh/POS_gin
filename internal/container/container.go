package container

import (
	"POS/internal/handlers"
	"POS/internal/repositories"
	"POS/internal/services"

	"gorm.io/gorm"
)

type Container struct {
	Auth     *handlers.AuthHandler
	Category *handlers.CategoryHandler
	Product  *handlers.ProductHandler
}

func NewContainer(db *gorm.DB) *Container {
	// ---------- Auth ----------
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// ---------- Category ----------
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo, categoryRepo)
	productHandler := handlers.NewProductHandler(productService)

	return &Container{
		Auth:     authHandler,
		Category: categoryHandler,
		Product:  productHandler,
	}
}
