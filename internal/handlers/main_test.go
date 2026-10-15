package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"

	"POS/internal/models"
	"POS/internal/repositories"

	"github.com/gin-gonic/gin"
)

// doJSON mengirim request JSON dan mengembalikan response recorder.
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

// uintToStr convert uint ke string untuk dipakai di URL.
func uintToStr(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

// seedCategory menyimpan kategori langsung ke mock repository.
func seedCategory(repo *repositories.MockCategoryRepository, name, desc string) *models.Category {
	ctx := context.Background()
	cat := &models.Category{Name: name, Description: desc}
	_ = repo.Create(ctx, cat)
	return cat
}
