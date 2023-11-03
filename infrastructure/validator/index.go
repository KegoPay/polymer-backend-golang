package validator

func init(){
	validate.RegisterValidation("password", validatePasswordStrength)
}

type Validator struct {}

func (v *Validator) ValidateStruct(payload interface{}) *[]error {
	return validateStruct(payload)
}

var ValidatorInstance = Validator{}