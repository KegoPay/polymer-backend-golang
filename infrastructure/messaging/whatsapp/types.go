package sms

type TermiiOTPResponse struct {
	PinID		string `json:"pinId"`
	Recipient	string `json:"to"`
	Msg			string `json:"smsStatus"`
}

type SMSServiceType interface {
	SendSMS(phone string, sms string) *string
}