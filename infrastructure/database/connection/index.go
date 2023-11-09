package connection

import (
	"kego.com/infrastructure/database/connection/cache"
	"kego.com/infrastructure/database/connection/datastore"
)

func ConnectToDatabase(){
	datastore.ConnectToDatabase()
	cache.ConnectToCache()
}