package container

import (
	"POS/internal/handlers"
	"POS/internal/repositories"
	"POS/internal/services"

	"gorm.io/gorm"
)

type Container struct {
	Auth        *handlers.AuthHandler
	Category    *handlers.CategoryHandler
	Product     *handlers.ProductHandler
	Stock       *handlers.StockHandler
	Transaction *handlers.TransactionHandler
	Setting     *handlers.SettingHandler
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

	// ---------- Product ----------
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo, categoryRepo)
	productHandler := handlers.NewProductHandler(productService)

	// ---------- Stock ----------
	stockRepo := repositories.NewStockRepository(db)
	stockService := services.NewStockService(stockRepo, productRepo)
	stockHandler := handlers.NewStockHandler(stockService)

	// ---------- Setting ----------
	settingRepo := repositories.NewSettingRepository(db)
	settingService := services.NewSettingService(settingRepo)
	settingHandler := handlers.NewSettingHandler(settingService)

	// ---------- Transaction ----------
	txnRepo := repositories.NewTransactionRepository(db)
	txnService := services.NewTransactionService(txnRepo, productRepo, settingService)
	txnHandler := handlers.NewTransactionHandler(txnService)

	return &Container{
		Auth:        authHandler,
		Category:    categoryHandler,
		Product:     productHandler,
		Stock:       stockHandler,
		Transaction: txnHandler,
		Setting:     settingHandler,
	}
}
