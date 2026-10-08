package validators_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"POS/internal/validators"

	"github.com/gin-gonic/gin"
)

type TestInput struct {
	Name  string `json:"name" binding:"required,min=3"`
	Age   int    `json:"age" binding:"required,gte=17"`
	Email string `json:"email" binding:"required,email"`
}

func doBind(t *testing.T, body string) map[string]string {
	t.Helper()
	gin.SetMode(gin.TestMode)

	var input TestInput
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	err := c.ShouldBindJSON(&input)
	if err == nil {
		t.Fatal("harusnya error")
	}
	return validators.FormatValidationError(err)
}

func TestFormatValidationError_FieldRules(t *testing.T) {
	body := `{"name":"Ab","age":10,"email":"bukan-email"}`
	result := doBind(t, body)

	// field pakai JSON tag (lowercase)
	if result["name"] == "" {
		t.Errorf("name harusnya ada pesan. Result: %+v", result)
	}
	if result["age"] == "" {
		t.Errorf("age harusnya ada pesan. Result: %+v", result)
	}
	if result["email"] == "" {
		t.Errorf("email harusnya ada pesan. Result: %+v", result)
	}

	if !strings.Contains(result["name"], "minimal 3") {
		t.Errorf("pesan name salah: %s", result["name"])
	}
	if !strings.Contains(result["email"], "tidak valid") {
		t.Errorf("pesan email salah: %s", result["email"])
	}
}

func TestFormatValidationError_JSONSyntax(t *testing.T) {
	body := `{"name":"Budi",`
	result := doBind(t, body)

	if result["body"] == "" {
		t.Errorf("harusnya ada pesan body. Result: %+v", result)
	}
}

func TestFormatValidationError_WrongType(t *testing.T) {
	body := `{"name":"Budi","age":"dua puluh","email":"a@b.com"}`
	result := doBind(t, body)

	if result["age"] == "" {
		t.Errorf("age harusnya ada pesan tipe salah. Result: %+v", result)
	}
}

func TestFormatValidationError_EmptyBody(t *testing.T) {
	body := ``
	result := doBind(t, body)

	if result["body"] == "" {
		t.Errorf("harusnya ada pesan body kosong. Result: %+v", result)
	}
}

func TestFormatValidationError_ValidJSON(t *testing.T) {
	body := `{}`
	result := doBind(t, body)

	if result["name"] == "" || result["age"] == "" || result["email"] == "" {
		t.Errorf("harusnya semua field required error: %+v", result)
	}
}

func TestFormatValidationError_NeverEmpty(t *testing.T) {
	cases := []string{
		``,
		`{}`,
		`{"name":"x"}`,
		`{"name":"Budi",`,
		`{"name":"Budi","age":"salah"}`,
		`not even json`,
	}

	for _, body := range cases {
		t.Run(body, func(t *testing.T) {
			result := doBind(t, body)
			if len(result) == 0 {
				t.Errorf("map kosong untuk body: %s", body)
			}
		})
	}
}

func TestFormatValidationError_Marshalable(t *testing.T) {
	body := `{"name":"Ab"}`
	result := doBind(t, body)

	_, err := json.Marshal(result)
	if err != nil {
		t.Errorf("harusnya bisa di-marshal: %v", err)
	}
}
