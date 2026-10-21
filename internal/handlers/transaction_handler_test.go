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

func setupTxnRouter() (*gin.Engine, *repositories.MockProductRepository, *repositories.MockCategoryRepository) {
	gin.SetMode(gin.TestMode)

	catRepo := repositories.NewMockCategoryRepository()
	prodRepo := repositories.NewMockProductRepository(catRepo)
	stockRepo := repositories.NewMockStockRepository(prodRepo)
	txnRepo := repositories.NewMockTransactionRepository(prodRepo, stockRepo)
	settingRepo := repositories.NewMockSettingRepository()
	_ = settingRepo.Upsert(context.Background(), "tax_rate", "11")
	_ = settingRepo.Upsert(context.Background(), "invoice_prefix", "INV")

	settingSvc := services.NewSettingService(settingRepo)
	txnSvc := services.NewTransactionService(txnRepo, prodRepo, settingSvc)
	h := handlers.NewTransactionHandler(txnSvc)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", float64(1))
		c.Set("role", "admin")
		c.Next()
	})

	r.POST("/transactions", h.Create)
	r.GET("/transactions", h.GetAll)
	r.GET("/transactions/:id", h.GetByID)
	r.POST("/transactions/:id/pay", h.Pay)
	r.POST("/transactions/:id/void", h.Void)

	return r, prodRepo, catRepo
}

func seedTxnProductHandler(t *testing.T, prodRepo *repositories.MockProductRepository, catRepo *repositories.MockCategoryRepository, sku, name string, price float64, stock int) uint {
	t.Helper()
	ctx := context.Background()
	cat := &models.Category{Name: "Test"}
	_ = catRepo.Create(ctx, cat)
	p := &models.Product{SKU: sku, Name: name, Price: price, Stock: stock, CategoryID: cat.ID, IsActive: true}
	_ = prodRepo.Create(ctx, p)
	return p.ID
}

func TestTransactionHandler_Create_Pending(t *testing.T) {
	r, prodRepo, catRepo := setupTxnRouter()
	pid := seedTxnProductHandler(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 100)

	body := map[string]interface{}{
		"items": []map[string]interface{}{
			{"product_id": pid, "quantity": 2},
		},
	}
	w := doJSON(r, "POST", "/transactions", body)

	if w.Code != http.StatusCreated {
		t.Errorf("harusnya 201, dapat: %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["status"] != "pending" {
		t.Errorf("status harusnya pending, dapat: %v", data["status"])
	}
}

func TestTransactionHandler_Create_Paid(t *testing.T) {
	r, prodRepo, catRepo := setupTxnRouter()
	pid := seedTxnProductHandler(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 100)

	body := map[string]interface{}{
		"items": []map[string]interface{}{
			{"product_id": pid, "quantity": 2},
		},
		"payment_method": "cash",
		"payment_amount": 10000,
	}
	w := doJSON(r, "POST", "/transactions", body)

	if w.Code != http.StatusCreated {
		t.Errorf("harusnya 201, dapat: %d", w.Code)
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["status"] != "paid" {
		t.Errorf("status harusnya paid, dapat: %v", data["status"])
	}
}

func TestTransactionHandler_Create_InsufficientStock(t *testing.T) {
	r, prodRepo, catRepo := setupTxnRouter()
	pid := seedTxnProductHandler(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 3)

	body := map[string]interface{}{
		"items": []map[string]interface{}{
			{"product_id": pid, "quantity": 10},
		},
	}
	w := doJSON(r, "POST", "/transactions", body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("harusnya 400, dapat: %d", w.Code)
	}
}

func TestTransactionHandler_Create_ValidationError(t *testing.T) {
	r, _, _ := setupTxnRouter()

	body := map[string]interface{}{
		"items": []map[string]interface{}{},
	}
	w := doJSON(r, "POST", "/transactions", body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("harusnya 400, dapat: %d", w.Code)
	}
}

func TestTransactionHandler_Pay(t *testing.T) {
	r, prodRepo, catRepo := setupTxnRouter()
	pid := seedTxnProductHandler(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 100)

	// Create pending
	w := doJSON(r, "POST", "/transactions", map[string]interface{}{
		"items": []map[string]interface{}{{"product_id": pid, "quantity": 2}},
	})
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	id := uint(data["id"].(float64))

	// Pay
	w = doJSON(r, "POST", "/transactions/"+uintToStr(id)+"/pay", map[string]interface{}{
		"payment_method": "cash",
		"payment_amount": 10000,
	})
	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestTransactionHandler_Void(t *testing.T) {
	r, prodRepo, catRepo := setupTxnRouter()
	pid := seedTxnProductHandler(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 100)

	// Create
	w := doJSON(r, "POST", "/transactions", map[string]interface{}{
		"items":          []map[string]interface{}{{"product_id": pid, "quantity": 2}},
		"payment_method": "cash",
		"payment_amount": 10000,
	})
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	id := uint(data["id"].(float64))

	// Void
	w = doJSON(r, "POST", "/transactions/"+uintToStr(id)+"/void", map[string]interface{}{
		"reason": "Salah input",
	})
	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d. Body: %s", w.Code, w.Body.String())
	}

	var resp2 map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp2)
	data2 := resp2["data"].(map[string]interface{})
	if data2["status"] != "void" {
		t.Errorf("status harusnya void, dapat: %v", data2["status"])
	}

	// Stok kembali
	p, _ := prodRepo.FindByID(context.Background(), pid)
	if p.Stock != 100 {
		t.Errorf("stok harusnya kembali 100, dapat: %d", p.Stock)
	}
}

func TestTransactionHandler_GetByID(t *testing.T) {
	r, prodRepo, catRepo := setupTxnRouter()
	pid := seedTxnProductHandler(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 100)

	w := doJSON(r, "POST", "/transactions", map[string]interface{}{
		"items": []map[string]interface{}{{"product_id": pid, "quantity": 1}},
	})
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	id := uint(data["id"].(float64))

	w = doJSON(r, "GET", "/transactions/"+uintToStr(id), nil)
	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d", w.Code)
	}
}
