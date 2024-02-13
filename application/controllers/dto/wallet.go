package dto

import "kego.com/entities"

type SendPaymentDTO struct {
	Pin         			string       	 `json:"pin"`
	FullName         		*string      	 `json:"fullName"`
	Amount      			uint64      	 `json:"amount"`
	DestinationCountryCode  string 		 	 `json:"destinationCountryCode"`
	BankCode				string 		 	 `json:"bankCode"`
	BranchCode				*string 		 `json:"branchCode"`
	AccountNumber 			string 			 `json:"accountNumber"`
	Description 			*string 		 `json:"description"`
	IPAddress 				string 			 `json:"ipAddress"`
}

type NameVerificationDTO struct {
	AccountNumber  string       `bson:"accountNumber" json:"accountNumber"`
	BankName       string       `bson:"bankName" json:"bankName"`
}

type SetPaymentTagDTO struct {
	Tag  string       `bson:"tag" json:"tag" validate:"required,alphanum,string_min_length_3"`
}

type UpdateAddressDTO struct {
	State 	  string       `bson:"state" json:"state" validate:"required,alpha,oneof=lagos"`
	LGA    	  string       `bson:"lga" json:"lga" validate:"required,alpha,mx=20"`
	Street	  string       `bson:"street" json:"street" validate:"required,alpha,max=300"`
}

type ToggleNotificationOptionsDTO struct {
	Emails 			 *bool `bson:"emails" json:"emails"`
	PushNotification *bool `bson:"pushNotification" json:"pushNotification"`
}

type EmailSubscriptionDTO struct {
	Email 	 string 					   `json:"email"`
	Channel	 entities.SubscriptionChannels `json:"channel"`
}

type FlutterwaveWebhookDTO struct {
	EventType	string						`json:"event.type"`
	Transfer 	FlutterwaveWebhookTransfer	`json:"transfer"`
}

type FlutterwaveWebhookTransfer struct {
	ID				uint							`json:"id"`
	Status			string							`json:"status"`
	Ref				string							`json:"reference"`
	Msg				string							`json:"complete_message"`
	Currency		string							`json:"currency"`
	RecepientName	string							`json:"fullname"`
	Amount			float32							`json:"amount"`
	Meta   			FlutterwaveWebhookTransferMeta  `json:"meta"`
}

type FlutterwaveWebhookTransferMeta struct {
	WalletID	string 	`json:"walletID"`
	UserID		string 	`json:"userID"`
}
