package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"POS/internal/handlers"
	"POS/internal/models"
	"POS/internal/repositories"
	"POS/internal/services"

	"github.com/gin-gonic/gin"
)

func setupStockRouter() (*gin.Engine, *repositories.MockProductRepository, *repositories.MockCategoryRepository) {
	gin.SetMode(gin.TestMode)

	catRepo := repositories.NewMockCategoryRepository()
	prodRepo := repositories.NewMockProductRepository(catRepo)
	stockRepo := repositories.NewMockStockRepository(prodRepo)
	svc := services.NewStockService(stockRepo, prodRepo)
	h := handlers.NewStockHandler(svc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", float64(1))
		c.Set("role", "admin")
		c.Next()
	})

	r.POST("/stock/mutations", h.Create)
	r.GET("/stock/mutations", h.GetAll)
	r.GET("/stock/mutations/:id", h.GetByID)
	r.GET("/products/:id/mutations", h.GetByProductID)

	return r, prodRepo, catRepo
}

func seedStockProduct(t *testing.T, prodRepo *repositories.MockProductRepository, catRepo *repositories.MockCategoryRepository, stock int) uint {
	t.Helper()
	ctx := context.Background()
	cat := &models.Category{Name: "Makanan"}
	_ = catRepo.Create(ctx, cat)

	prod := &models.Product{
		SKU: "MKN-001", Name: "Indomie", Price: 3500,
		Stock: stock, CategoryID: cat.ID, IsActive: true,
	}
	_ = prodRepo.Create(ctx, prod)
	return prod.ID
}

func TestStockHandler_Create_In_Success(t *testing.T) {
	r, prodRepo, catRepo := setupStockRouter()
	prodID := seedStockProduct(t, prodRepo, catRepo, 10)

	body := map[string]interface{}{
		"product_id": prodID, "type": "in", "quantity": 50, "note": "Restock",
	}
	w := doJSON(r, "POST", "/stock/mutations", body)

	if w.Code != http.StatusCreated {
		t.Errorf("harusnya 201, dapat: %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["stock_after"].(float64) != 60 {
		t.Errorf("stock_after harusnya 60, dapat: %v", data["stock_after"])
	}
}

func TestStockHandler_Create_Out_Insufficient(t *testing.T) {
	r, prodRepo, catRepo := setupStockRouter()
	prodID := seedStockProduct(t, prodRepo, catRepo, 5)

	body := map[string]interface{}{
		"product_id": prodID, "type": "out", "quantity": 50,
	}
	w := doJSON(r, "POST", "/stock/mutations", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("harusnya 400, dapat: %d", w.Code)
	}
}

func TestStockHandler_GetByProductID(t *testing.T) {
	r, prodRepo, catRepo := setupStockRouter()
	prodID := seedStockProduct(t, prodRepo, catRepo, 10)

	_ = doJSON(r, "POST", "/stock/mutations", map[string]interface{}{
		"product_id": prodID, "type": "in", "quantity": 20,
	})
	_ = doJSON(r, "POST", "/stock/mutations", map[string]interface{}{
		"product_id": prodID, "type": "out", "quantity": 5,
	})

	w := doJSON(r, "GET", "/products/"+uintToStr(prodID)+"/mutations", nil)
	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d", w.Code)
	}
}
