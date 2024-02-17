package entities

type Country struct{
	Name 			 string					`json:"name"`
	ISOCode 		 string					`json:"isoCode"`
	FlagURL			 string					`json:"flagURL"`
	Symbol			 string					`json:"symbol"`
	Rate			 string					`json:"rate"`
	ServicesAllowed  []CountryServiceType	`json:"servicesAllowed"`
}