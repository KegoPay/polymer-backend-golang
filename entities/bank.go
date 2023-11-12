package entities

type Bank struct{
	Name 		 string		`json:"name" validate:"required"`
	Code 		 string		`json:"code" validate:"required"`
	LongCode 	 string		`json:"longCode"`
}