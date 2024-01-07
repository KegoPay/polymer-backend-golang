package validator

import (
	"os"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"kego.com/entities"
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

func userAgentConditionalValidator(fl validator.FieldLevel) bool {
	agent := fl.Field().String()
	if os.Getenv("GIN_MODE") != "release" {
		return true
	}
	if !strings.Contains(agent, "Android") ||  !strings.Contains(agent, "iOS") {
		return false
	}
	return true
}

func exclusiveEmailAndPhone(fl validator.FieldLevel) bool {
	email, emailOK := fl.Field().Interface().(*string)
	phone, phoneOK := fl.Parent().FieldByName("Phone").Interface().(*entities.PhoneNumber)
	if emailOK && phoneOK && ((email != nil && phone == nil) || (email == nil && phone != nil)) {
		return true
	}
	return false
}

func stringLengthValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return len(value) >= 3
}