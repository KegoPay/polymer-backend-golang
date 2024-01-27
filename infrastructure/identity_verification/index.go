package identityverification

import (
	"os"

	prembly_identity_verification "kego.com/infrastructure/identity_verification/prembly"
	identity_verification_types "kego.com/infrastructure/identity_verification/types"
	"kego.com/infrastructure/network"
)

var IdentityVerifier identity_verification_types.IdentityVerifierType

func InitialiseIdentityVerifier(){
	IdentityVerifier = &prembly_identity_verification.PremblyIdentityVerification{
		Network: &network.NetworkController{
			BaseUrl: os.Getenv("PREMBLY_BASE_URL"),
		},
		API_KEY: os.Getenv("PREMBLY_API_KEY"),
		APP_ID: os.Getenv("PREMBLY_APP_ID"),
	}
}