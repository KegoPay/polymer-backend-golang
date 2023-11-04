package entities

// Represents the services supported for each country
type CountryServiceType string

var (
	SignUp CountryServiceType = 		 "signup"
	InstantTransfer CountryServiceType = "instant_transfer"
	MobileMoney CountryServiceType = 	 "mobile_money"
)