package entities

import "time"

// Represents the services supported for each country
type CountryServiceType string

var (
	SignUp CountryServiceType = 		 "signup"
	InstantTransfer CountryServiceType = "instant_transfer"
	MobileMoney CountryServiceType = 	 "mobile_money"
)

// Represents the reason an account was restricted
type AccountRestrictedReason struct {
	Reason string `json:",omitempty"`
	Duration *time.Time
}