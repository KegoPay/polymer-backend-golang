package entities

type Country struct{
	Name 			 string					`json:"name"`
	ISOCode 		 string					`json:"isoCode"`
	FlagURL			 string					`json:"flagURL"`
	ServicesAllowed  []CountryServiceType	`json:"servicesAllowed"`
}