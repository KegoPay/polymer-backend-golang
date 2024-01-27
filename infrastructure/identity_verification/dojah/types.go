package dojah_identity_verification

import identity_verification_types "kego.com/infrastructure/identity_verification/types"

type DojahBVNResponse struct {
	Data        	identity_verification_types.BVNData `json:"entity"`
}
