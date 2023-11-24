package validator

func init(){
	validate.RegisterValidation("exclusive_email_phone", exclusiveEmailAndPhone, true)
	validate.RegisterValidation("password", validatePasswordStrength)
	validate.RegisterValidation("user_agent", userAgentConditionalValidator)
}

type Validator struct {}

func (v *Validator) ValidateStruct(payload interface{}) *[]error {
	return validateStruct(payload)
}

var ValidatorInstance = Validator{}