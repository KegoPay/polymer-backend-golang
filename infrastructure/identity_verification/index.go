package identityverification

import (
	"os"

	"kego.com/infrastructure/network"
)

var IdentityVerifier PaystackVerification

func InitialiseIdentityVerifier(){
	IdentityVerifier = PaystackVerification{
		Network: &network.NetworkController{
			BaseUrl: os.Getenv("PAYSTACK_BASE_URL"),
		},
		AuthToken: os.Getenv("PAYSTACK_ACCESS_TOKEN"),
	}
}