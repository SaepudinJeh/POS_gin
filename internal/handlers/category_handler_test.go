package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"POS/internal/handlers"
	"POS/internal/repositories"
	"POS/internal/services"

	"github.com/gin-gonic/gin"
)

func setupCategoryRouter() (*gin.Engine, *repositories.MockCategoryRepository) {
	gin.SetMode(gin.TestMode)

	mockRepo := repositories.NewMockCategoryRepository()
	svc := services.NewCategoryService(mockRepo)
	h := handlers.NewCategoryHandler(svc)

	r := gin.New()
	r.POST("/categories", h.Create)
	r.GET("/categories", h.GetAll)
	r.GET("/categories/:id", h.GetByID)
	r.PUT("/categories/:id", h.Update)
	r.DELETE("/categories/:id", h.Delete)

	return r, mockRepo
}

func TestCategoryHandler_Create_Success(t *testing.T) {
	r, _ := setupCategoryRouter()

	body := map[string]string{"name": "Makanan", "description": "Kategori makanan"}
	w := doJSON(r, "POST", "/categories", body)

	if w.Code != http.StatusCreated {
		t.Errorf("harusnya 201, dapat: %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestCategoryHandler_Create_ValidationError(t *testing.T) {
	r, _ := setupCategoryRouter()

	body := map[string]string{"description": "Tanpa nama"}
	w := doJSON(r, "POST", "/categories", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("harusnya 400, dapat: %d", w.Code)
	}
}

func TestCategoryHandler_GetAll(t *testing.T) {
	r, _ := setupCategoryRouter()

	_ = doJSON(r, "POST", "/categories", map[string]string{"name": "Makanan"})
	_ = doJSON(r, "POST", "/categories", map[string]string{"name": "Minuman"})

	w := doJSON(r, "GET", "/categories", nil)

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

func TestCategoryHandler_GetByID_Success(t *testing.T) {
	r, mockRepo := setupCategoryRouter()
	cat := seedCategory(mockRepo, "Makanan", "Kategori makanan")

	w := doJSON(r, "GET", "/categories/"+uintToStr(cat.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d", w.Code)
	}
}

func TestCategoryHandler_GetByID_NotFound(t *testing.T) {
	r, _ := setupCategoryRouter()

	w := doJSON(r, "GET", "/categories/999", nil)
	if w.Code != http.StatusNotFound {
		t.Errorf("harusnya 404, dapat: %d", w.Code)
	}
}

func TestCategoryHandler_Update_Success(t *testing.T) {
	r, mockRepo := setupCategoryRouter()
	cat := seedCategory(mockRepo, "Makanan", "Kategori makanan")

	body := map[string]string{"name": "Minuman", "description": "Kategori minuman"}
	w := doJSON(r, "PUT", "/categories/"+uintToStr(cat.ID), body)

	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestCategoryHandler_Delete_Success(t *testing.T) {
	r, mockRepo := setupCategoryRouter()
	cat := seedCategory(mockRepo, "Makanan", "Kategori makanan")

	w := doJSON(r, "DELETE", "/categories/"+uintToStr(cat.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d", w.Code)
	}
}

func TestCategoryHandler_Delete_NotFound(t *testing.T) {
	r, _ := setupCategoryRouter()

	w := doJSON(r, "DELETE", "/categories/999", nil)
	if w.Code != http.StatusNotFound {
		t.Errorf("harusnya 404, dapat: %d", w.Code)
	}
}
