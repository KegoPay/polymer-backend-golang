package entities

import (
	"time"

	"kego.com/application/utils"
)

type TransactionIntent string

const (
	PaystackDVACredit          TransactionIntent = "paystack_dva_credit"
	FlutterwaveDVACredit          TransactionIntent = "flutterwave_dva_credit"
	ChimoneyDebitInternational TransactionIntent = "chimoney_debit_international"
	PaystackDebitLocal         TransactionIntent = "paystack_debit_local"
	FlutterwaveDebitLocal         TransactionIntent = "flutterwave_debit_local"
)

type DeviceInfo struct {
	IPAddress  string    `bson:"ipAddress" json:"ipAddress" validate:"required,ip"`
	DeviceID   string    `bson:"deviceID" json:"deviceID" validate:"required"`
	UserAgent  string 	 `bson:"userAgent" json:"userAgent" validate:"user_agent,required"`
}

type TransactionSender struct {
	BusinessName string      `bson:"businessName" json:"businessName" validate:"required"`
	FirstName    string      `bson:"firstName" json:"firstName" validate:"required"`
	LastName     string      `bson:"lastName" json:"lastName" validate:"required"`
	Email        string      `bson:"email" json:"email,omitempty" validate:"required,email"`
}

type TransactionRecepient struct {
	Name     	  string  `bson:"name" json:"name" validate:"required"`
	AccountNumber string  `bson:"accountNumber" json:"accountNumber" validate:"required"`
	Country       string  `bson:"country" json:"country" validate:"iso3166_1_alpha2"`
	BankCode      string  `bson:"bankCode" json:"bankCode" validate:"required"`
	BankName      string  `bson:"bankName" json:"bankName"`
	BranchCode    *string `bson:"branchCode" json:"branchCode" validate:"required"`
}

type Transaction struct {
	TransactionReference string               `bson:"transactionReference" json:"transactionReference" validate:"required"`
	Amount               uint64               `bson:"amount" json:"amount" validate:"required"`
	AmountInNGN          uint64          	  `bson:"amountInNGN" json:"amountInNGN" validate:"required"`
	Fee          		 uint64          	  `bson:"fee" json:"fee" validate:"required"`
	ProcessorFee         uint64          	  `bson:"processorFee" json:"processorFee" validate:"required"`
	AmountInUSD          *uint64              `bson:"amountInUSD" json:"amountInUSD" validate:"required"`
	Currency             string               `bson:"currency" json:"currency" validate:"iso4217"`
	ProcessorFeeCurrency string               `bson:"processorFeeCurrency" json:"processorFeeCurrency" validate:"iso4217"`
	WalletID             string               `bson:"walletID" json:"walletID" validate:"required"`
	UserID               string               `bson:"userID" json:"userID" validate:"required"`
	BusinessID           *string              `bson:"businessID" json:"businessID"`
	Description          string               `bson:"description" json:"description" validate:"required"`
	MetaData          	 any               	  `bson:"metadata" json:"metadata" validate:"required"`
	Location             Location          	  `bson:"location" json:"location" validate:"required"`
	Intent               TransactionIntent 	  `bson:"intent" json:"intent" validate:"required"`
	DeviceInfo           DeviceInfo        	  `bson:"deviceInfo" json:"deviceInfo" validate:"required"`
	Sender               TransactionSender 	  `bson:"transactionSender" json:"transactionSender" validate:"required"`
	Recepient            TransactionRecepient `bson:"transactionRcepient" json:"transactionRcepient" validate:"required"`

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

