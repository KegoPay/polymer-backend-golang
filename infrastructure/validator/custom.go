package validator

import (
	"fmt"
	"regexp"
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
		}else {
			return false
		}
	}
	return digitCount == 6
}

func validateTrxPinStrength(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	digitCount := 0

	for _, char := range password {
		if unicode.IsDigit(char) {
			digitCount++
		}
		return false
	}
	return digitCount == 4
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

func alphaSpace(fl validator.FieldLevel) bool {
	pattern := `^[a-zA-Z ]+$`

	// Compile the regular expression
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regular expression:", err)
		return false
	}

	// Match the string against the regular expression
	return re.MatchString(fl.Field().String())
}