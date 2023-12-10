package dto

type SendPaymentDTO struct {
	Pin         string       `bson:"pin" json:"pin"`
	Amount      uint64       `bson:"amount" json:"amount"`
}