package biometric

import (
	"os"

	faceapi "usepolymer.co/infrastructure/biometric/faceAPI"
	"usepolymer.co/infrastructure/network"
)

var BiometricService BiometricServiceType

func InitialiseBiometricService() {
	BiometricService = &faceapi.FaceAPIBiometricService{
		Network: &network.NetworkController{
			BaseUrl: os.Getenv("FACEAPI_BASE_URL"),
		},
		API_KEY: os.Getenv("FACEAPI_API_KEY"),
	}
}
