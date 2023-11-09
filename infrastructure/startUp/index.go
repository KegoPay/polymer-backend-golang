package startup

import (
	"kego.com/infrastructure/database"
	"kego.com/infrastructure/database/connection/datastore"
	"kego.com/infrastructure/logger"
)

// Used to start services such as loggers, databases, queues, etc.
func StartServices(){
	// initialise logger module
	logger.InitializeLogger()
	// set up databases
	database.SetUpDatabase()
}

// Used to clean up after services that have been shutdown.
func CleanUpServices(){
	// clean up database resources
	datastore.CleanUp()
}