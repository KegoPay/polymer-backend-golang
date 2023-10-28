package startup

import "kego.com/infrastructure/logger"

// Used to start services such as loggers, databases, queues, etc.
func StartServices(){
	// initialise logger module
	logger.InitializeLogger()
}

// Used to clean up after services that have been shutdown.
func CleanUpServices(){}