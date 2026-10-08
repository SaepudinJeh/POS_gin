package validators

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {
	result := make(map[string]string)

	// 1. Error dari validator (field rules)
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, e := range validationErrs {
			result[e.Field()] = messageForTag(e)
		}
		return result
	}

	// 2. JSON syntax error
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		result["body"] = fmt.Sprintf("Format JSON tidak valid (posisi byte %d)", syntaxErr.Offset)
		return result
	}

	// 3. Tipe data salah
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		field := typeErr.Field
		if field == "" {
			field = "body"
		}
		result[field] = fmt.Sprintf("Tipe data salah, harusnya %s", typeErr.Type.String())
		return result
	}

	// 4. Body kosong
	if errors.Is(err, io.EOF) {
		result["body"] = "Body request tidak boleh kosong"
		return result
	}

	// 5. Error lain (fallback)
	msg := err.Error()
	if strings.Contains(msg, "Unsupported Media Type") {
		result["header"] = "Content-Type harus application/json"
		return result
	}

	result["body"] = msg
	return result
}

func messageForTag(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()
	param := e.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s wajib diisi", field)
	case "email":
		return fmt.Sprintf("Format %s tidak valid", field)
	case "min":
		return fmt.Sprintf("%s minimal %s karakter", field, param)
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", field, param)
	case "gte":
		return fmt.Sprintf("%s harus lebih besar atau sama dengan %s", field, param)
	case "gt":
		return fmt.Sprintf("%s harus lebih besar dari %s", field, param)
	case "lte":
		return fmt.Sprintf("%s harus lebih kecil atau sama dengan %s", field, param)
	case "lt":
		return fmt.Sprintf("%s harus lebih kecil dari %s", field, param)
	case "notspace":
		return fmt.Sprintf("%s tidak boleh mengandung spasi", field)
	case "numeric":
		return fmt.Sprintf("%s harus berupa angka", field)
	case "alphanum":
		return fmt.Sprintf("%s hanya boleh huruf dan angka", field)
	default:
		return fmt.Sprintf("%s tidak valid", field)
	}
}
