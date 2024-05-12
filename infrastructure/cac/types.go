package cac_service

type CACBusinessStatus string
var ACTIVE CACBusinessStatus = "ACTIVE"
var INACTIVE CACBusinessStatus = "INACTIVE"

type CACBusiness struct {
	RCNumber			string				`json:"rcNumber" validate:"required,max=30"`
	Status				CACBusinessStatus	`json:"companyStatus" validate:"oneof=ACTIVE INACTIVE"`
	FullAddress			string				`json:"address" validate:"required,max=300"`
	RegistrationDate	string				`json:"registrationDate" validate:"required,max=50"`
	Name				string				`json:"approvedName" validate:"required,max=100"`
	Email				*string				`json:"email" validate:"max=100"`
}

type CACBusinessNameSearchResponse struct {
	Data	*[]CACBusiness	`json:"data"`
	Message	string	`json:"message"`
	Error	uint	`json:"errorCode"`
}