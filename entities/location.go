package entities

type Location struct	{
	IPAddress  	 string `bson:"ipAddress" json:"ipAddress" validate:"latitude"`
}
