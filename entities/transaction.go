package entities

import "time"

type TransactionIntent string

const (
	PaystackDVACredit          TransactionIntent = "paystack_dva_credit"
	ChimoneyDebitInternational TransactionIntent = "chimoney_debit_international"
	PaystackDebitLocal         TransactionIntent = "paystack_debit_local"
)

type DeviceInfo struct {
	IPAddress  string    `bson:"ipAddress" json:"ipAddress" validate:"required,ip"`
	DeviceID   string    `bson:"deviceID" json:"deviceID" validate:"required"`
	UserAgent  UserAgent `bson:"userAgent" json:"userAgent" validate:"user_agent,required"`
	AppVersion string    `bson:"appVersion" json:"appVersion" validate:"required"`
}

type TransactionSender struct {
	UserID       string      `bson:"userID" json:"userID" validate:"required"`
	BusinessID   string      `bson:"businessID" json:"businessID" validate:"required"`
	BusinessName string      `bson:"businessName" json:"businessName" validate:"required"`
	FirstName    string      `bson:"firstName" json:"firstName" validate:"required"`
	LastName     string      `bson:"lastName" json:"lastName" validate:"required"`
	Email        string      `bson:"email" json:"email,omitempty" validate:"required,email"`
	Phone        PhoneNumber `bson:"phone" json:"phone,omitempty" validate:"required"`
}

type TransactionRecepient struct {
	FirstName     string  `bson:"firstName" json:"firstName" validate:"required"`
	LastName      string  `bson:"lastName" json:"lastName" validate:"required"`
	BankName      string  `bson:"bankName" json:"bankName" validate:"required,email"`
	AccountNumber string  `bson:"accountNumber" json:"accountNumber" validate:"required"`
	Country       string  `bson:"country" json:"country" validate:"iso3166_1_alpha2"`
	BankCode      string  `bson:"bankCode" json:"bankCode" validate:"required"`
	BranchCode    *string `bson:"branchCode" json:"branchCode" validate:"required"`
}

type Transaction struct {
	TransactionReference string            `bson:"transactionReference" json:"transactionReference" validate:"required"`
	Amount               uint64            `bson:"amount" json:"amount" validate:"required"`
	Currency             string            `bson:"currency" json:"currency" validate:"iso4217"`
	WalletID             uint64            `bson:"walletID" json:"walletID" validate:"required"`
	UserID               uint64            `bson:"userID" json:"userID" validate:"required"`
	BusinessID           string            `bson:"businessID" json:"businessID" validate:"required"`
	Description          string            `bson:"description" json:"description" validate:"required"`
	Location             Location          `bson:"location" json:"location" validate:"required"`
	Intent               TransactionIntent `bson:"intent" json:"intent" validate:"required"`
	DeviceInfo           DeviceInfo        `bson:"deviceInfo" json:"deviceInfo" validate:"required"`
	Sender               TransactionSender `bson:"transactionSender" json:"transactionSender" validate:"required"`

	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
