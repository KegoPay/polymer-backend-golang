package identity_verification_types

type IdentityVerifierType interface {
	FetchBVNDetails(string) (*BVNData, error)
	EmailVerification(email string) (bool, error)
}


type BVNData struct {
	Gender            string    `json:"gender"`
	WatchListed       string    `json:"watch_listed"`
	FirstName         string    `json:"first_name"`
	MiddleName        *string    `json:"middle_name"`
	LastName          string    `json:"last_name"`
	DateOfBirth       string    `json:"date_of_birth"`
	PhoneNumber       string     `json:"phone_number1"`
	Nationality       string    `json:"nationality"`
	Base64Image       string    `json:"image"`
}