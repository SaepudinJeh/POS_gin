package validators

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// FormatValidationError mengubah error dari validator menjadi map yang rapi.
func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			tag := fieldError.Tag()
			param := fieldError.Param()

			var message string

			switch tag {
			case "required":
				message = fmt.Sprintf("%s wajib diisi", field)
			case "email":
				message = fmt.Sprintf("Format %s tidak valid", field)
			case "min":
				message = fmt.Sprintf("%s minimal %s karakter", field, param)
			case "max":
				message = fmt.Sprintf("%s maksimal %s karakter", field, param)
			case "gt":
				message = fmt.Sprintf("%s harus lebih besar dari %s", field, param)
			case "gte":
				message = fmt.Sprintf("%s harus lebih besar atau sama dengan %s", field, param)
			case "lt":
				message = fmt.Sprintf("%s harus lebih kecil dari %s", field, param)
			case "lte":
				message = fmt.Sprintf("%s harus lebih kecil atau sama dengan %s", field, param)
			case "notspace":
				message = fmt.Sprintf("%s tidak boleh mengandung spasi", field)
			default:
				message = fmt.Sprintf("%s tidak valid", field)
			}

			errors[field] = message
		}
	}

	return errors
}
