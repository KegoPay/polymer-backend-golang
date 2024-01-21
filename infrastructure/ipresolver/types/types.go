package types

type IPResolver interface {
	ConnectToDB()
	LookUp(ipAddress string) (*IPResult, error)
}

type IPResult struct {
	AcuracyRadius 			 	int 			 `bson:"acuracyRadius" json:"acuracyRadius"`
	Longitude 					float64 			 `bson:"longitude" json:"longitude"`
	Latitude  					float64				 `bson:"latitude" json:"latitude"`
	City						string 	 			 `bson:"city" json:"city"`
	CountryCode		  			string 				 `bson:"countryCode" json:"countryCode"`
}