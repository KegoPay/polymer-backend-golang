package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()


func validateStruct(payload interface{}) *[]error {
	err := validate.Struct(payload)
	if err == nil {
		return nil
	}
	var ve validator.ValidationErrors
	errs := []error{}
	if errors.As(err, &ve) {
		for _, fe := range ve {
			errs = append(errs, errors.New(msgForTag(fe)))
		}
	}
	return &errs
}

func msgForTag(fe validator.FieldError) string {
	fe.Tag()
	err_msg := fieldErrorMap(fe.Tag(), fe.Field(), fe.Value(), fe.Param())
	return err_msg
}
