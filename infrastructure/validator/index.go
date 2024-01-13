package validator

func init(){
	validate.RegisterValidation("exclusive_email_phone", exclusiveEmailAndPhone, true)
	validate.RegisterValidation("password", validatePasswordStrength)
	validate.RegisterValidation("trx_pin", validateTrxPinStrength)
	validate.RegisterValidation("user_agent", userAgentConditionalValidator)
	validate.RegisterValidation("string_min_length_3", stringLengthValidator)
}

type Validator struct {}

func (v *Validator) ValidateStruct(payload interface{}) *[]error {
	return validateStruct(payload)
}

var ValidatorInstance = Validator{}