package dto



type FlutterwaveWebhookDTO struct {
	EventType			*string						`json:"event.type"`
	TrxRef 				*string						`json:"txRef"`
	Amount 				*string						`json:"amount"`
	ChargedAmount 		*string						`json:"charged_amount"`
	Status 				*string						`json:"status"`
	IPAddress 			*string						`json:"ip"`
	Currency 			*string						`json:"currency"`
	AppFee 				*string						`json:"appfee"`
	Transfer 			*FlutterwaveWebhookTransfer	`json:"transfer"`
	Entity 				*FlutterwaveWebhookEntity	`json:"entity"`
	Customer 			*FlutterwaveWebhookCustomer	`json:"Customer"`
}

type FlutterwaveWebhookEntity struct {
	AccNumber 				string						`json:"account_number"`
	FirstName 				string						`json:"first_name"`
	LastName 				string						`json:"last_name"`
}

type FlutterwaveWebhookCustomer struct {
	FullName 			string						`json:"fullName"`
	Email 				string						`json:"email"`
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
