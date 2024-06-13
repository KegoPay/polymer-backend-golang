package types

type LocalPaymentProcessorType interface {
	InitialisePaymentProcessor()
	NameVerification(accountNumber string, bankCode string) (*NameVerificationResponseField, *int, error)
	InitiateMobileMoneyTransfer(payload *InitiateLocalTransferPayload) (*InitiateLocalTransferDataField, *int, error)
	InitiateLocalTransfer(payload *InitiateLocalTransferPayload) (*InitiateLocalTransferDataField, *int, error)
	GenerateDVA(payload *any) (*VirtualAccountPayload, *int, error)
}

type NameVerificationResponseDTO struct {
	Status  bool                          `json:"status"`
	Message string                        `json:"message"`
	Data    NameVerificationResponseField `json:"data"`
}

type NameVerificationResponseField struct {
	AccountName   string `json:"account_name"`
	AccountNumber int    `json:"account_number"`
}

type InitiateLocalTransferPayload struct {
	Reference   string                   `json:"reference"`
	Destination LocalTransferDestination `json:"destination"`
}

type LocalTransferDestination struct {
	Type        string                              `json:"type"`
	Amount      float32                              `json:"amount"`
	Currency    string                              `json:"currency"`
	Narration   string                              `json:"narration"`
	BankAccount LocalTransferDestinationBankAccount `json:"bank_account"`
	Customer    LocalTransferDestinationCustomer    `json:"customer"`
}

type LocalTransferDestinationBankAccount struct {
	Bank      string `json:"bank"`
	Account   string `json:"account"`
	Currency  string `json:"currency"`
	Narration string `json:"narration"`
}

type LocalTransferDestinationCustomer struct {
	Name string `json:"name"`
}

type InitiateTransferPayloadResponse struct {
	Status  string                         `json:"status"`
	Message string                         `json:"message"`
	Data    InitiateLocalTransferDataField `json:"data"`
}

type InitiateLocalTransferDataField struct {
	BankName  string  `json:"bank_name"`
	FullName  string  `json:"full_name"`
	Fee       float32 `json:"fee"`
	BankCode  string  `json:"bank_code"`
	PaymentID uint    `json:"id"`
}

type CreateVirtualAccountPayload struct {
	Email                string  `json:"email"`
	Permanent            bool    `json:"is_permanent"`
	BVN                  string  `json:"bvn"`
	TransactionReference string  `json:"tx_ref"`
	FirstName            string  `json:"firstname"`
	LastName             string  `json:"lastname"`
	Narration            string  `json:"narration"`
	Amount               *uint64 `json:"amount"`
	Currency             string  `json:"currency"`
}

type CreateVirtualAccountResponse struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Data    VirtualAccountPayload `json:"data"`
}

type VirtualAccountPayload struct {
	AccoutName    string  `json:"account_name"`
	AccountNumber string  `json:"account_number"`
	BankCode      string  `json:"bank_code"`
	BankName      string  `json:"bank_name"`
	ID            string  `json:"unique_id"`
	CreatedAt     float32 `json:"created_at"`
}
