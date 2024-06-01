package prembly_identity_verification

import identity_verification_types "usepolymer.co/infrastructure/identity_verification/types"

type PremblyBVNResponse struct {
	Status  bool                                `json:"status"`
	Detail  string                              `json:"detail"`
	Message string                              `json:"message"`
	Data    identity_verification_types.BVNData `json:"data"`
}

type PremblyFaceMatchResponse struct {
	Status     bool    `json:"status"`
	Detail     string  `json:"detail"`
	Message    string  `json:"message"`
	Confidence float32 `json:"confidence"`
}
