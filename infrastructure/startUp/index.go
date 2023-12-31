package startup

import (
	"kego.com/infrastructure/database"
	"kego.com/infrastructure/database/connection/datastore"
	fileupload "kego.com/infrastructure/file_upload"
	identityverification "kego.com/infrastructure/identity_verification"
	"kego.com/infrastructure/logger"
	international_payment_processor "kego.com/infrastructure/payment_processor/chimoney"
	local_payment_processor "kego.com/infrastructure/payment_processor/paystack"
)

// Used to start services such as loggers, databases, queues, etc.
func StartServices(){
	// initialise logger module
	logger.InitializeLogger()
	// set up databases
	database.SetUpDatabase()
	
	fileupload.InitialiseFileUploader()
	identityverification.InitialiseIdentityVerifier()
	local_payment_processor.InitialisePaystackPaymentProcessor()
	international_payment_processor.InitialiseChimoneyPaymentProcessor()
}

// Used to clean up after services that have been shutdown.
func CleanUpServices(){
	// clean up database resources
	datastore.CleanUp()
}