package biometric

import (
	"os"

	"kego.com/infrastructure/biometric/prembly"
	"kego.com/infrastructure/network"
)

var BiometricService BiometricServiceType

func InitialiseBiometricService(){
	BiometricService = &prembly.PremblyBiometricService{
		Network: &network.NetworkController{
			BaseUrl: os.Getenv("PREMBLY_BASE_URL"),
		},
		API_KEY: os.Getenv("PREMBLY_API_KEY"),
		APP_ID: os.Getenv("PREMBLY_APP_ID"),
	}
}