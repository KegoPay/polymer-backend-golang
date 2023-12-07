package entities

type Location struct	{
	Longitude 	 string `bson:"longitude" json:"longitude" validate:"longitude"` // Two-letter country code (ISO 3166-1 alpha-2)
	Latitude  string `bson:"latitude" json:"latitude" validate:"latitude"`
}
