package identity_verification_types

type IdentityVerifierType interface {
	FetchBVNDetails(string) (*BVNData, error)
	FetchNINDetails(string) (*NINData, error)
	EmailVerification(email string) (bool, error)
	FetchAdvancedCACDetails(rc string) (*CompanyProfile, error)
	ImgLivenessCheck(img string) (bool, error)
}

type BVNData struct {
	Gender        string  `json:"gender"`
	WatchListed   string  `json:"watch_listed"`
	FirstName     string  `json:"first_name"`
	MiddleName    *string `json:"middle_name"`
	LastName      string  `json:"last_name"`
	DateOfBirth   string  `json:"date_of_birth"`
	PhoneNumber   string  `json:"phone_number1"`
	Nationality   string  `json:"nationality"`
	Address       string  `json:"residential_address"`
	Base64Image   string  `json:"image"`
	NIN           string  `json:"nin"`
	LGAOfOrigin   string  `json:"lga_of_origin"`
	StateOfOrigin string  `json:"state_of_origin"`
	Title         string  `json:"title"`
}

type NINData struct {
	Gender      string  `json:"gender"`
	FirstName   string  `json:"first_name"`
	MiddleName  *string `json:"middle_name"`
	LastName    string  `json:"last_name"`
	DateOfBirth string  `json:"date_of_birth"`
	PhoneNumber *string `json:"phone_number"`
	Address     string  `json:"address"`
	Base64Image string  `json:"photo"`
}

type CACResponse struct {
	Data CompanyProfile `json:"entity"`
}

type CompanyProfile struct {
	City       string             `json:"City"`
	LGA        string             `json:"LGA"`
	State      string             `json:"State"`
	Status     string             `json:"Status"`
	Affiliates []AffiliateProfile `json:"affiliates"`
}

type LivenessCheckResult struct {
	Entity LivenessCheckResultEntity `json:"entity"`
}

type LivenessCheckResultEntity struct {
	Liveness LivenessCheckResultLiveness `json:"liveness"`
}

type LivenessCheckResultLiveness struct {
	LivenessCheck       bool    `json:"liveness_check"`
	LivenessProbability float32 `json:"liveness_probability"`
}

type AffiliateProfile struct {
	Name          string  `json:"name"`
	Status        string  `json:"status"`
	IDType        string  `json:"identityType"`
	IDNumber      string  `json:"identityNumber"`
	AffiliateType string  `json:"affiliateType"`
	ShareAllotted string  `json:"shareAllotted"`
	Email         *string `json:"email"`
	Phone         *string `json:"phone"`
}
