package validator

import "fmt"

func fieldErrorMap(tag string, field string, value interface{}, param interface{}) string {
	err_map := map[string]string{
		"required": 				fmt.Sprintf("%s is required", field),
		"excludes": 				fmt.Sprintf(`"%s" is not allowed in %s`, value, field),
		"min":      				fmt.Sprintf("%s cannot be less than %s digits", field, param),
		"max":      				fmt.Sprintf("%s cannot be more than %s digits", field, param),
		"email":      				fmt.Sprintf("%s is not a valid email", value),
		"iso3166_1_alpha2": 		fmt.Sprintf("%s should be a 2 letter country code (ISO 3166-1 alpha-2)", field),
		"oneof": 					fmt.Sprintf("%s must be one of %s", field, param),
		// custom
		"password":      			fmt.Sprintf("%s should be a secret 4 digit number", field),
		"exclusive_email_phone": 	"An email or phone number must be provided to sign up",
		
	}
	return err_map[tag]
}
