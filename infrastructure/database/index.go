package database

import (
	"usepolymer.co/infrastructure/database/connection"
)

func SetUpDatabase() {
	connection.ConnectToDatabase()
}

type BaseModel interface {
	ParseModel() any
}
