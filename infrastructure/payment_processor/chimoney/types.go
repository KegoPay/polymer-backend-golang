package chimoney_international_payment_processor

import "kego.com/entities"


type ChimoneyExchangeRateDTO struct {
	Error 	string	   			    `json:"error"`
	Status 	string					`json:"status"`
	Data	entities.ExchangeRates	`json:"data"`
	ValidTill	 int64				`json:"expiresAtTimestamp"`
}

type ChimoneySupportedBanksDTO struct {
	Error 	string	   `json:"error"`
	Status 	string		    `json:"status"`
	Data    []entities.Bank `json:"data"`
}

type InternationalPaymentRequestResponsePayload struct {
	Error 	string	   `json:"error"`
	Status 	string	   `json:"status"`
	Data    InternationalPaymentRequestResponseDataPayload     `json:"data"`
}

type InternationalPaymentRequestPayload struct {
	DestinationCountry string `json:"countryToSend"`
	BankCode string `json:"account_bank"`
	AccountNumber string `json:"account_number"`
	Reference string `json:"reference"`
	ValueInUSD float32 `json:"valueInUSD"`
	FullName *string `json:"fullname"`
}

type InternationalPaymentRequestResponseDataPayload struct {
	Chimoneys []InternationalPaymentRequestResponseDataChimoneyPayload `json:"chimoneys"`
}

type InternationalPaymentRequestResponseDataChimoneyPayload struct {
	ID string `json:"id"`
	AccountNumber string `json:"account_number"`
	BankCode string `json:"account_bank"`
	Type string `json:"type"`
	ChiRef string `json:"chiRef"`
	Status string `json:"status"`
	ValueInUSD float32 `json:"valueInUSD"`
	CountrySentTo string `json:"countryToSend"`
	Fee float32 `json:"fee"`
	IssueDate string `json:"issueDate"`
	Issuer string `json:"issuer"`
	IssueID string `json:"issueID"`
}
