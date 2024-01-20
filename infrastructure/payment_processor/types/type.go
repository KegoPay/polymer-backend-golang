package types

type LocalPaymentProcessorType interface {
	InitialisePaymentProcessor()
	NameVerification(accountNumber string, bankCode string) (*NameVerificationResponseField, *int, error)
	InitiateLocalTransfer(payload *InitiateLocalTransferPayload) (*InitiateLocalTransferDataField, *int, error)
	GenerateDVA(payload *CreateVirtualAccountPayload) (*VirtualAccountPayload, *int, error)
}

type NameVerificationResponseDTO struct {
	Status 	bool						  `json:"status"`
	Message string						  `json:"message"`
	Data	NameVerificationResponseField `json:"data"`
}

type NameVerificationResponseField struct {
	AccountName 	string  `json:"account_name"`
	AccountNumber  	int		`json:"account_number"`
}


type InitiateLocalTransferPayload struct {
	AccountBank 	string		`json:"account_bank"`
	AccountNumber 	string		`json:"account_number"`
	Amount 			uint64		`json:"amount"`
	Narration 		string		`json:"narration"`
	Currency 		string		`json:"currency"`
	Reference 		string		`json:"reference"`
	CallbackURL 	string		`json:"callback_url"`
	DebitCurrency 	string		`json:"debit_currency"`
}

type InitiateLocalTransferPayloadResponse struct {
	Status 		string									`json:"status"`
	Message 	string									`json:"message"`
	Data		InitiateLocalTransferDataField			`json:"data"`
}

type InitiateLocalTransferDataField struct {
	BankName 	string			`json:"bank_name"`
	FullName 	string			`json:"full_name"`
	Fee			float32			`json:"fee"`
}

type CreateVirtualAccountPayload struct {
	Email 					string		`json:"email"`
	Permanent 				bool		`json:"is_permanent"`
	BVN 					string		`json:"bvn"`
	TransactionReference 	string		`json:"tx_ref"`
	FirstName 				string		`json:"firstname"`
	LastName 				string		`json:"lastname"`
	Narration 				string		`json:"narration"`
	Currency 				string		`json:"currency"`
}

type CreateVirtualAccountResponse struct {
	Status 		string					`json:"status"`
	Message 	string					`json:"message"`
	Data		VirtualAccountPayload	`json:"data"`
}

type VirtualAccountPayload struct {
	FlutterwaveReference 	string							`json:"flw_ref"`
	OrderReference			string							`json:"order_ref"`
	AccountNumber		 	string							`json:"account_number"`
	AccountStatus		 	string							`json:"account_status"`
	BankName 				string							`json:"bank_name"`
	Status 					string							`json:"status"`
	Message 				string							`json:"message"`
	Amount 					float32							`json:"amount"`
	Note 					string							`json:"note"`
	CreatedAt 				float32							`json:"created_at"`
	ExpiryDate 				float32							`json:"1703031769350"`
}
