package entities

type Location struct	{
	IPAddress  	 string `bson:"ipAddress" json:"ipAddress" validate:"ip"`
	Latitude  	 float64 `bson:"latitude" json:"latitude" validate:"latitude"`
	Longitude  	 float64 `bson:"longitude" json:"longitude" validate:"longitude"`
}
