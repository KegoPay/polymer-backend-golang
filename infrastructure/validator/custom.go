package validator

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

func validatePasswordStrength(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	digitCount := 0

	for _, char := range password {
		if unicode.IsDigit(char) {
			digitCount++
		}
	}
	return digitCount >= 4
}