package dto

type SendPaymentDTO struct {
	Pin         string       `bson:"pin" json:"pin"`
	Amount      uint64       `bson:"amount" json:"amount"`
}

type NameVerificationDTO struct {
	AccountNumber  string       `bson:"accountNumber" json:"accountNumber"`
	BankName       string       `bson:"bankName" json:"bankName"`
}