package connection

import (
	"usepolymer.co/infrastructure/database/connection/cache"
	"usepolymer.co/infrastructure/database/connection/datastore"
)

func ConnectToDatabase() {
	datastore.ConnectToDatabase()
	cache.ConnectToCache()
}
