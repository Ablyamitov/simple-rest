package validation

import (
	"fmt"
	"github.com/Ablyamitov/simple-rest/internal/app/wrapper"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

var validate = validator.New()

func Validate(obj interface{}) error {

	if err := validate.RegisterValidation("notblank", NotBlank); err != nil {
		wrapper.LogError(fmt.Sprintf("Error register notblank validation: %v", err),
			"validation.Validate")
	}
	return validate.Struct(obj)
}

func NotBlank(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(field.String())) > 0
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}
