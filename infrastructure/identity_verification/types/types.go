package identity_verification_types

type IdentityVerifierType interface {
	FetchBVNDetails(string) (*BVNData, error)
	FetchNINDetails(string) (*NINData, error)
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
	Address	 	      string    `json:"residential_address"`
	Base64Image       string    `json:"image"`
}

type NINData struct {
	Gender            string    `json:"gender"`
	FirstName         string    `json:"firstname"`
	MiddleName        *string    `json:"middlename"`
	LastName          string    `json:"surname"`
	DateOfBirth       string    `json:"birthdate"`
	PhoneNumber       *string     `json:"telephoneno"`
	Nationality       string    `json:"birth_country"`
	Address	 	      string    `json:"address"`
	Base64Image       string    `json:"photo"`
}