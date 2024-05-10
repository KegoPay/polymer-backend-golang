package dto

type FlutterwaveWebhookDTO struct {
    EventType      string                         `form:"event.type"`
    TrxRef         *string                         `form:"txRef"`
    Amount         *float32                         `form:"amount"`
    ChargedAmount  *string                         `form:"charged_amount"`
    Status         *string                         `form:"status"`
    IPAddress      *string                         `form:"IP"`
    Currency       *string                         `form:"currency"`
    AppFee         *string                         `form:"appfee"`
    MerchantFee    *string                         `form:"merchantfee"`
	Transfer 	  *FlutterwaveWebhookTransfer	`json:"transfer"`
    MerchantBearsFee *string                       `form:"merchantbearsfee"`
    Customer       *FlutterwaveWebhookCustomer     `form:"customer"`
    Entity         *FlutterwaveWebhookEntity       `form:"entity"`
}

type FlutterwaveWebhookEntity struct {
    AccNumber  string `form:"entity[account_number]"`
    FirstName  string `form:"entity[first_name]"`
    LastName   string `form:"entity[last_name]"`
}

type FlutterwaveWebhookCustomer struct {
    Phone       string `form:"customer[phone]"`
    FullName    string `form:"customer[fullName]"`
    Email       string `form:"customer[email]"`
}

type FlutterwaveWebhookTransfer struct {
	ID				uint							`json:"id"`
	Status			string							`json:"status"`
	Ref				string							`json:"reference"`
	Msg				string							`json:"complete_message"`
	Currency		string							`json:"currency"`
	RecepientName	string							`json:"fullname"`
	Amount			float32							`json:"amount"`
	Meta			FlutterwaveWebhookTransferMeta  `json:"meta"`
}


type FlutterwaveWebhookTransferMeta struct {
	WalletID	string 	`json:"walletID"`
	UserID		string 	`json:"userID"`
}

type ChimoneyWebhookDTO struct {
    
}