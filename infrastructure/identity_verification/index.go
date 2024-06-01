package identityverification

import (
	"os"

	dojah_identity_verification "usepolymer.co/infrastructure/identity_verification/dojah"
	identity_verification_types "usepolymer.co/infrastructure/identity_verification/types"
	"usepolymer.co/infrastructure/network"
)

var IdentityVerifier identity_verification_types.IdentityVerifierType

func InitialiseIdentityVerifier() {
	IdentityVerifier = &dojah_identity_verification.DojahIdentityVerification{
		Network: &network.NetworkController{
			BaseUrl: os.Getenv("DOJAH_BASE_URL"),
		},
		API_KEY: os.Getenv("DOJAH_API_KEY"),
		APP_ID:  os.Getenv("DOJAH_APP_ID"),
	}
}
