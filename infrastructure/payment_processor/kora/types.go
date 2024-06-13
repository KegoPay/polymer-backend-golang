package kora_local_payment_processor

type CreateVirtualAccountPayload struct {
	Reference string      `json:"account_reference"`
	Permanent bool        `json:"permanent"`
	Name      string      `json:"account_name"`
	BankCode  string      `json:"bank_code"`
	Customer  DVACustomer `json:"customer"`
	KYC       DVAKYC      `json:"kyc"`
}

type DVACustomer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DVAKYC struct {
	NIN string `json:"nin"`
	BVN string `json:"bvn"`
}
