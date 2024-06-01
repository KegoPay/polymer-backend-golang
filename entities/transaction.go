package entities

import (
	"time"

	"usepolymer.co/application/utils"
)

type TransactionIntent string
type TransactionStatus string

const (
	PaystackDVACredit     TransactionIntent = "paystack_dva_credit"
	FlutterwaveDVACredit  TransactionIntent = "flutterwave_dva_credit"
	InternationalDebit    TransactionIntent = "international_debit"
	InternationalDebitFee TransactionIntent = "internation_debit_fee"
	PaystackDebitLocal    TransactionIntent = "paystack_debit_local"
	LocalDebit            TransactionIntent = "local_debit"
	LocalDebitFee         TransactionIntent = "local_debit_fee"
	LocalDebitVAT         TransactionIntent = "local_debit_vat"
	PolymerVAT            TransactionIntent = "polymer_vat"
	PolymerFee            TransactionIntent = "polymer_fee"

	TransactionPending   TransactionStatus = "pending"
	TransactionCompleted TransactionStatus = "completed"
	TransactionFailed    TransactionStatus = "failed"
)

type DeviceInfo struct {
	IPAddress string  `bson:"ipAddress" json:"ipAddress" validate:"required,ip"`
	DeviceID  *string `bson:"deviceID" json:"deviceID" validate:"required"`
	UserAgent *string `bson:"userAgent" json:"userAgent" validate:"user_agent,required"`
}

type TransactionSender struct {
	BusinessName  *string `bson:"businessName" json:"businessName" validate:"required"`
	FullName      string  `bson:"fullName" json:"fullName" validate:"required"`
	Email         *string `bson:"email" json:"email,omitempty" validate:"required,email"`
	BankCode      *string `bson:"bankCode" json:"bankCode" validate:"required"`
	BankName      *string `bson:"bankName" json:"bankName"`
	BranchCode    *string `bson:"branchCode" json:"branchCode"`
	AccountNumber *string `bson:"accountNumber" json:"accountNumber" validate:"required"`
}

type TransactionRecepient struct {
	FullName      string  `bson:"fullName" json:"fullName" validate:"required"`
	AccountNumber string  `bson:"accountNumber" json:"accountNumber" validate:"required"`
	Country       *string `bson:"country" json:"country" validate:"iso3166_1_alpha2"`
	BankCode      *string `bson:"bankCode" json:"bankCode" validate:"required"`
	BankName      *string `bson:"bankName" json:"bankName"`
	BranchCode    *string `bson:"branchCode" json:"branchCode" validate:"required"`
}

type Transaction struct {
	TransactionReference string               `bson:"transactionReference" json:"transactionReference" validate:"required"`
	Amount               uint64               `bson:"amount" json:"amount" validate:"required"`
	AmountInNGN          *uint64              `bson:"amountInNGN" json:"amountInNGN" validate:"required"`
	TotalAmountInNGN     *uint64              `bson:"totalAmountInNGN" json:"totalAmountInNGN" validate:"required"`
	Fee                  uint64               `bson:"fee" json:"fee" validate:"required"`
	Vat                  uint64               `bson:"vat" json:"vat" validate:"required"`
	ProcessorFee         uint64               `bson:"processorFee" json:"processorFee" validate:"required"`
	ProcessorFeeVAT      uint64               `bson:"processorFeeVAT" json:"processorFeeVAT" validate:"required"`
	AmountInUSD          *uint64              `bson:"amountInUSD" json:"amountInUSD" validate:"required"`
	Currency             string               `bson:"currency" json:"currency" validate:"iso4217"`
	ProcessorFeeCurrency string               `bson:"processorFeeCurrency" json:"processorFeeCurrency" validate:"iso4217"`
	WalletID             string               `bson:"walletID" json:"walletID" validate:"required"`
	UserID               string               `bson:"userID" json:"userID" validate:"required"`
	BusinessID           *string              `bson:"businessID" json:"businessID"`
	Message              *string              `bson:"message" json:"message"`
	Description          string               `bson:"description" json:"description" validate:"required"`
	Status               TransactionStatus    `bson:"status" json:"status" validate:"required"`
	MetaData             any                  `bson:"metadata" json:"metadata" validate:"required"`
	Location             Location             `bson:"location" json:"location" validate:"required"`
	Intent               TransactionIntent    `bson:"intent" json:"intent" validate:"required"`
	DeviceInfo           *DeviceInfo          `bson:"deviceInfo" json:"deviceInfo"`
	Sender               TransactionSender    `bson:"transactionSender" json:"transactionSender" validate:"required"`
	Recepient            TransactionRecepient `bson:"transactionRecepient" json:"transactionRecepient" validate:"required"`

	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

func (trx Transaction) ParseModel() any {
	if trx.ID == "" {
		trx.CreatedAt = time.Now()
		trx.ID = utils.GenerateUUIDString()
	}
	trx.UpdatedAt = time.Now()
	return &trx
}
