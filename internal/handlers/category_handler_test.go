package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"POS/internal/handlers"
	"POS/internal/models"
	"POS/internal/repositories"
	"POS/internal/services"

	"github.com/gin-gonic/gin"
)

// setupRouter bikin router minimal untuk test handler kategori.
// Kita tidak pakai auth middleware di sini — fokus test handler.
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

// helper: kirim request JSON dan kembalikan response recorder
func doJSON(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}

	req, _ := http.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ---------- CREATE ----------

func TestCategoryHandler_Create_Success(t *testing.T) {
	r, _ := setupCategoryRouter()

	body := map[string]string{
		"name":        "Makanan",
		"description": "Kategori makanan",
	}
	w := doJSON(r, "POST", "/categories", body)

	if w.Code != http.StatusCreated {
		t.Errorf("harusnya 201, dapat: %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["success"] != true {
		t.Errorf("success harusnya true")
	}
}

func TestCategoryHandler_Create_ValidationError(t *testing.T) {
	r, _ := setupCategoryRouter()

	// name kosong -> gagal validasi required
	body := map[string]string{
		"description": "Tanpa nama",
	}
	w := doJSON(r, "POST", "/categories", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("harusnya 400, dapat: %d", w.Code)
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["success"] != false {
		t.Errorf("success harusnya false")
	}
	if resp["errors"] == nil {
		t.Errorf("harusnya ada field 'errors'")
	}
}

func TestCategoryHandler_Create_DuplicateName(t *testing.T) {
	r, _ := setupCategoryRouter()

	body := map[string]string{
		"name":        "Makanan",
		"description": "Kategori makanan",
	}
	_ = doJSON(r, "POST", "/categories", body)

	// Kirim yang sama lagi
	w := doJSON(r, "POST", "/categories", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("harusnya 400, dapat: %d", w.Code)
	}
}

// ---------- GET ALL ----------

func TestCategoryHandler_GetAll_Empty(t *testing.T) {
	r, _ := setupCategoryRouter()

	w := doJSON(r, "GET", "/categories", nil)

	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d", w.Code)
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	data, _ := resp["data"].([]interface{})
	if len(data) != 0 {
		t.Errorf("harusnya kosong, dapat: %d", len(data))
	}
}

func TestCategoryHandler_GetAll_WithData(t *testing.T) {
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

// ---------- GET BY ID ----------

func TestCategoryHandler_GetByID_Success(t *testing.T) {
	r, mockRepo := setupCategoryRouter()

	// Seed langsung ke mock supaya dapat ID yang pasti
	cat := seedCategory(mockRepo, "Makanan", "Kategori makanan")

	w := doJSON(r, "GET", "/categories/"+uintToStr(cat.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestCategoryHandler_GetByID_NotFound(t *testing.T) {
	r, _ := setupCategoryRouter()

	w := doJSON(r, "GET", "/categories/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("harusnya 404, dapat: %d", w.Code)
	}
}

func TestCategoryHandler_GetByID_InvalidID(t *testing.T) {
	r, _ := setupCategoryRouter()

	w := doJSON(r, "GET", "/categories/abc", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("harusnya 400, dapat: %d", w.Code)
	}
}

// ---------- UPDATE ----------

func TestCategoryHandler_Update_Success(t *testing.T) {
	r, mockRepo := setupCategoryRouter()
	cat := seedCategory(mockRepo, "Makanan", "Kategori makanan")

	body := map[string]string{
		"name":        "Minuman",
		"description": "Kategori minuman",
	}
	w := doJSON(r, "PUT", "/categories/"+uintToStr(cat.ID), body)

	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestCategoryHandler_Update_ValidationError(t *testing.T) {
	r, mockRepo := setupCategoryRouter()
	cat := seedCategory(mockRepo, "Makanan", "Kategori makanan")

	body := map[string]string{"name": ""} // name kosong
	w := doJSON(r, "PUT", "/categories/"+uintToStr(cat.ID), body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("harusnya 400, dapat: %d", w.Code)
	}
}

func TestCategoryHandler_Update_NotFound(t *testing.T) {
	r, _ := setupCategoryRouter()

	body := map[string]string{"name": "Minuman"}
	w := doJSON(r, "PUT", "/categories/999", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("harusnya 400, dapat: %d", w.Code)
	}
}

// ---------- DELETE ----------

func TestCategoryHandler_Delete_Success(t *testing.T) {
	r, mockRepo := setupCategoryRouter()
	cat := seedCategory(mockRepo, "Makanan", "Kategori makanan")

	w := doJSON(r, "DELETE", "/categories/"+uintToStr(cat.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("harusnya 200, dapat: %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestCategoryHandler_Delete_NotFound(t *testing.T) {
	r, _ := setupCategoryRouter()

	w := doJSON(r, "DELETE", "/categories/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("harusnya 404, dapat: %d", w.Code)
	}
}

// ---------- helper ----------

func seedCategory(repo *repositories.MockCategoryRepository, name, desc string) *models.Category {
	cat := &models.Category{Name: name, Description: desc}
	_ = repo.Create(cat)
	return cat
}

func uintToStr(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}
