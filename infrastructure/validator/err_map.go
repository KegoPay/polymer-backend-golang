package validator

import "fmt"

func fieldErrorMap(tag string, field string, value interface{}, param interface{}) string {
	err_map := map[string]string{
		"required": 				fmt.Sprintf("%s is required", field),
		"excludes": 				fmt.Sprintf(`"%s" is not allowed in %s`, value, field),
		"min":      				fmt.Sprintf("%s cannot be less than %s digits", field, param),
		"max":      				fmt.Sprintf("%s cannot be more than %s digits", field, param),
		"email":      				fmt.Sprintf("%s is not a valid email", value),
		"password":      			fmt.Sprintf("%s should be a secret 4 digit number", field),
		"iso3166_1_alpha3": 		fmt.Sprintf("%s should be a 3 letter country code (ISO 3166-1 alpha-3)", field),
	}
	return err_map[tag]
}
