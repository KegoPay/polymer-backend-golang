package startup

import (
	"kego.com/infrastructure/database/connection/cache"
	"kego.com/infrastructure/database/connection/datastore"
	"kego.com/infrastructure/logger"
)

// Used to start services such as loggers, databases, queues, etc.
func StartServices(){
	// initialise logger module
	logger.InitializeLogger()
	// connect to database
	datastore.ConnectToDatabase()
	// connect to cache
	cache.ConnectToCache()
}

// Used to clean up after services that have been shutdown.
func CleanUpServices(){
	// clean up database resources
	datastore.CleanUp()
}