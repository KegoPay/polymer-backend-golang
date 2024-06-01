package dojah_identity_verification

import identity_verification_types "usepolymer.co/infrastructure/identity_verification/types"

type DojahBVNResponse struct {
	Data identity_verification_types.BVNData `json:"entity"`
}

type DojahNINResponse struct {
	Data  identity_verification_types.NINData `json:"entity"`
	Error string                              `json:"error"`
}

type DojahEmailVerification struct {
	Entity DojahEmailVerificationPayload `json:"entity"`
}

type DojahEmailVerificationPayload struct {
	Score         uint               `json:"score"`
	Deliverable   bool               `json:"deliverable"`
	DomainDetails DojahDomainDetails `json:"domain_details"`
}

type DojahDomainDetails struct {
	SusTLD     bool `json:"suspicious_tld"`
	Registered bool `json:"registered"`
	Disposable bool `json:"disposable"`
}
