package identity_verification_types

type IdentityVerifierType interface {
	FetchBVNDetails(string) (*BVNData, error)
	FaceMatch(string, string) (*float32, error)
}


type BVNData struct {
	Gender            string    `json:"gender"`
	WatchListed       string    `json:"watchListed"`
	FirstName         string    `json:"firstName"`
	MiddleName        *string    `json:"middleName"`
	LastName          string    `json:"lastName"`
	DateOfBirth       string    `json:"dateOfBirth"`
	PhoneNumber       string     `json:"phoneNumber1"`
	Nationality       string    `json:"nationality"`
	Base64Image       string    `json:"base64Image"`
}