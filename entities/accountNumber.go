package entities

type BankDetails struct	{
	Number 	 	string `bson:"number" json:"number" validate:"required"`
	BankName  		string `bson:"bankName" json:"bankName" validate:"required"`
}
