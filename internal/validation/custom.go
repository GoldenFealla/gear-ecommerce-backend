package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/goldenfealla/gear-manager/domain"
)

func ValidateIsGear(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if _, ok := domain.GearTypeMap[val]; ok {
		return true
	}
	return false
}
