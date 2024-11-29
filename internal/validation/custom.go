package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/goldenfealla/gear-manager/domain"
)

func ValidateIsGear(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	key := strings.ToLower(val)
	if _, ok := domain.GearTypeMap[key]; ok {
		return true
	}
	return false
}
