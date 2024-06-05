package biometric

import (
	"os"

	prembly_idpass "usepolymer.co/infrastructure/biometric/prembly"
	"usepolymer.co/infrastructure/biometric/types"
	"usepolymer.co/infrastructure/network"
)

var BiometricService types.BiometricServiceType

func InitialiseBiometricService() {
	BiometricService = &prembly_idpass.PremblyBiometricService{
		Network: &network.NetworkController{
			BaseUrl: os.Getenv("PREMBLY_BASE_URL"),
		},
		API_KEY: os.Getenv("PREMBLY_API_KEY"),
		APP_ID:  os.Getenv("PREMBLY_APP_ID"),
	}
}
