package countriessupported

import "kego.com/entities"

var countriesSupported =  []entities.Country{
	{
		Name:           "Nigeria",
		ISOCode:        "NG",
		FlagURL:        "upload.wikimedia.org/wikipedia/commons/7/79/Flag_of_Nigeria.svg",
		ServicesAllowed: []entities.CountryServiceType{entities.SignUp, entities.InstantTransfer},
	},
	{
		Name:           "Canada",
		ISOCode:        "CA",
		FlagURL:        "upload.wikimedia.org/wikipedia/commons/d/d9/Flag_of_Canada_%28Pantone%29.svg",
		ServicesAllowed: []entities.CountryServiceType{entities.InstantTransfer},
	},
	{
		Name:           "Ghana",
		ISOCode:        "GH",
		FlagURL:        "upload.wikimedia.org/wikipedia/commons/1/19/Flag_of_Ghana.svg",
		ServicesAllowed: []entities.CountryServiceType{entities.InstantTransfer, entities.MobileMoney},
	},
	{
		Name:           "India",
		ISOCode:        "IN",
		FlagURL:        "upload.wikimedia.org/wikipedia/en/4/41/Flag_of_India.svg",
		ServicesAllowed: []entities.CountryServiceType{entities.InstantTransfer},
	},
	{
		Name:           "Kenya",
		ISOCode:        "KE",
		FlagURL:        "upload.wikimedia.org/wikipedia/commons/4/49/Flag_of_Kenya.svg",
		ServicesAllowed: []entities.CountryServiceType{entities.InstantTransfer, entities.MobileMoney},
	},
	{
		Name:           "Mexico",
		ISOCode:        "MX",
		FlagURL:        "upload.wikimedia.org/wikipedia/commons/f/fc/Flag_of_Mexico.svg",
		ServicesAllowed: []entities.CountryServiceType{entities.InstantTransfer},
	},
	{
		Name:           "Rwanda",
		ISOCode:        "RW",
		FlagURL:        "upload.wikimedia.org/wikipedia/commons/1/17/Flag_of_Rwanda.svg",
		ServicesAllowed: []entities.CountryServiceType{entities.InstantTransfer, entities.MobileMoney},
	},
	{
		Name:           "South Africa",
		ISOCode:        "ZA",
		FlagURL:        "upload.wikimedia.org/wikipedia/commons/a/af/Flag_of_South_Africa.svg",
		ServicesAllowed: []entities.CountryServiceType{entities.InstantTransfer, entities.MobileMoney},
	},
	{
		Name:           "United Kingdom",
		ISOCode:        "GB",
		FlagURL:        "upload.wikimedia.org/wikipedia/en/a/ae/Flag_of_the_United_Kingdom.svg",
		ServicesAllowed: []entities.CountryServiceType{entities.InstantTransfer},
	},
	{
		Name:           "United States of America",
		ISOCode:        "US",
		FlagURL:        "upload.wikimedia.org/wikipedia/en/a/a4/Flag_of_the_United_States.svg",
		ServicesAllowed: []entities.CountryServiceType{entities.InstantTransfer},
	},
}

func FilterCountries(filter entities.CountryServiceType) []entities.Country {
	var selectedCountries = []entities.Country{}
	for _, c := range countriesSupported {
		for _, s := range c.ServicesAllowed {
			if s == filter{
				selectedCountries = append(selectedCountries, c)
				break
			}
		}
	}
	return selectedCountries
}