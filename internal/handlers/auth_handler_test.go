package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"POS/internal/handlers"
	"POS/internal/repositories"
	"POS/internal/services"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	mockRepo := repositories.NewMockUserRepository()
	svc := services.NewAuthService(mockRepo)
	h := handlers.NewAuthHandler(svc)

	r := gin.Default()
	r.POST("/register", h.Register)
	return r
}

func TestAuthHandler_Register_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-min-32-characters-long!!")
	r := setupRouter()

	body := map[string]string{
		"name":     "Budi",
		"email":    "budi@test.com",
		"password": "rahasia123",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("harusnya 201, dapat: %d", w.Code)
	}
}
