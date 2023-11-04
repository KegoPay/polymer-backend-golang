package entities

type Country struct{
	Name 			 string
	ISOCode 		 string
	FlagURL			 string
	ServicesAllowed  []CountryServiceType
}