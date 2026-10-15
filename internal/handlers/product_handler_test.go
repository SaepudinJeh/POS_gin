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

func setupProductRouter() (*gin.Engine, *repositories.MockCategoryRepository, *repositories.MockProductRepository) {
	gin.SetMode(gin.TestMode)

	catRepo := repositories.NewMockCategoryRepository()
	prodRepo := repositories.NewMockProductRepository(catRepo)
	prodSvc := services.NewProductService(prodRepo, catRepo)
	h := handlers.NewProductHandler(prodSvc)

	r := gin.New()
	r.POST("/products", h.Create)
	r.GET("/products", h.GetAll)
	r.GET("/products/:id", h.GetByID)
	r.PUT("/products/:id", h.Update)
	r.DELETE("/products/:id", h.Delete)

	return r, catRepo, prodRepo
}

func seedProductCategory(catRepo *repositories.MockCategoryRepository, name string) uint {
	ctx := context.Background()
	cat := &models.Category{Name: name}
	_ = catRepo.Create(ctx, cat)
	return cat.ID
}

func TestProductHandler_Create_Success(t *testing.T) {
	r, catRepo, _ := setupProductRouter()
	catID := seedProductCategory(catRepo, "Makanan")

	body := map[string]interface{}{
		"sku": "MKN-001", "name": "Indomie Goreng",
		"price": 3500, "cost": 2500, "stock": 100, "category_id": catID,
	}
	w := doJSON(r, "POST", "/products", body)

	if w.Code != http.StatusCreated {
		t.Errorf("harusnya 201, dapat: %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_Create_ValidationError(t *testing.T) {
	r, _, _ := setupProductRouter()

	body := map[string]interface{}{"description": "tanpa nama"}
	w := doJSON(r, "POST", "/products", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("harusnya 400, dapat: %d", w.Code)
	}
}

func TestProductHandler_Create_CategoryNotFound(t *testing.T) {
	r, _, _ := setupProductRouter()

	body := map[string]interface{}{
		"sku": "MKN-001", "name": "Indomie", "price": 3500, "category_id": 999,
	}
	w := doJSON(r, "POST", "/products", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("harusnya 400, dapat: %d", w.Code)
	}
}

func TestProductHandler_GetAll(t *testing.T) {
	ctx := context.Background()
	r, catRepo, prodRepo := setupProductRouter()
	catID := seedProductCategory(catRepo, "Makanan")

	prodSvc := services.NewProductService(prodRepo, catRepo)
	_, _ = prodSvc.Create(ctx, services.ProductInput{SKU: "MKN-001", Name: "Indomie", CategoryID: catID})
	_, _ = prodSvc.Create(ctx, services.ProductInput{SKU: "MKN-002", Name: "Mie Sedaap", CategoryID: catID})

	w := doJSON(r, "GET", "/products", nil)

	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d", w.Code)
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data, _ := resp["data"].([]interface{})
	if len(data) != 2 {
		t.Errorf("harusnya 2, dapat: %d", len(data))
	}
}

func TestProductHandler_GetByID_Success(t *testing.T) {
	ctx := context.Background()
	r, catRepo, prodRepo := setupProductRouter()
	catID := seedProductCategory(catRepo, "Makanan")

	prodSvc := services.NewProductService(prodRepo, catRepo)
	p, _ := prodSvc.Create(ctx, services.ProductInput{SKU: "MKN-001", Name: "Indomie", CategoryID: catID})

	w := doJSON(r, "GET", "/products/"+uintToStr(p.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d", w.Code)
	}
}

func TestProductHandler_GetByID_NotFound(t *testing.T) {
	r, _, _ := setupProductRouter()

	w := doJSON(r, "GET", "/products/999", nil)
	if w.Code != http.StatusNotFound {
		t.Errorf("harusnya 404, dapat: %d", w.Code)
	}
}

func TestProductHandler_Update_Success(t *testing.T) {
	ctx := context.Background()
	r, catRepo, prodRepo := setupProductRouter()
	catID := seedProductCategory(catRepo, "Makanan")

	prodSvc := services.NewProductService(prodRepo, catRepo)
	p, _ := prodSvc.Create(ctx, services.ProductInput{SKU: "MKN-001", Name: "Indomie", CategoryID: catID})

	body := map[string]interface{}{
		"sku": "MKN-001", "name": "Indomie Spesial", "price": 4000, "category_id": catID,
	}
	w := doJSON(r, "PUT", "/products/"+uintToStr(p.ID), body)

	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_Delete_Success(t *testing.T) {
	ctx := context.Background()
	r, catRepo, prodRepo := setupProductRouter()
	catID := seedProductCategory(catRepo, "Makanan")

	prodSvc := services.NewProductService(prodRepo, catRepo)
	p, _ := prodSvc.Create(ctx, services.ProductInput{SKU: "MKN-001", Name: "Indomie", CategoryID: catID})

	w := doJSON(r, "DELETE", "/products/"+uintToStr(p.ID), nil)
	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d", w.Code)
	}
}

func TestProductHandler_Delete_NotFound(t *testing.T) {
	r, _, _ := setupProductRouter()

	w := doJSON(r, "DELETE", "/products/999", nil)
	if w.Code != http.StatusNotFound {
		t.Errorf("harusnya 404, dapat: %d", w.Code)
	}
}
