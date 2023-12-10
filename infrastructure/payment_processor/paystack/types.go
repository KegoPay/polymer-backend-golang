package payment_processor

type PaystackNameVerificationResponseDTO struct {
	Status 	bool						  `json:"status"`
	Message string						  `json:"message"`
	Data	NameVerificationResponseField `json:"data"`
}

type NameVerificationResponseField struct {
	AccountName 	string  `json:"account_name"`
	AccountNumber  	int		`json:"account_number"`
}