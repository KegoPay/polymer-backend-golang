package sms

import (
	"os"

	"usepolymer.co/infrastructure/network"
)

var SMSService SMSServiceType

func InitSMSService() {
	SMSService = &TermiiService{
		Network: &network.NetworkController{
			BaseUrl: os.Getenv("TERMII_URL"),
		},
		API_KEY: os.Getenv("TERMII_API_KEY"),
	}
}
