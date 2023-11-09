package database

import (
	"kego.com/infrastructure/database/connection"
)

func SetUpDatabase(){
	connection.ConnectToDatabase()
}

type BaseModel interface {
	ParseModel() any
}
