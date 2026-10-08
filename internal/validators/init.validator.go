package validators

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 1. Pakai nama JSON sebagai field name di error
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" || name == "" {
				return fld.Name
			}
			return name
		})

		// 2. Custom validator: notspace (tidak boleh ada spasi)
		_ = v.RegisterValidation("notspace", func(fl validator.FieldLevel) bool {
			return !strings.Contains(fl.Field().String(), " ")
		})
	}
}
