package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Custom validation error message for tag
//
// See list of tag here: [Fields]
//
// [Fields]: https://pkg.go.dev/github.com/go-playground/validator/v10#readme-fields
func messageForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "gte":
		return fmt.Sprintf("The length in this field need to be greater than or equal %v", fe.Param())
	case "lte":
		return fmt.Sprintf("The length in this field need to be lesser than or equal %v", fe.Param())
	default:
		return fmt.Sprintf("Field '%s': '%v' must satisfy '%s' '%v' criteria", fe.Field(), fe.Value(), fe.Tag(), fe.Param())
	}
}

func GetValidationError(list validator.ValidationErrors) []*ValidationError {
	l := make([]*ValidationError, len(list))

	for i, v := range list {
		ve := &ValidationError{}

		ve.Field = v.Field()
		ve.Message = messageForTag(v)

		l[i] = ve
	}

	return l
}
