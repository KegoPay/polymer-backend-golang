package startup

import (
	"kego.com/application"
	"kego.com/infrastructure/database"
	"kego.com/infrastructure/database/connection/datastore"
	fileupload "kego.com/infrastructure/file_upload"
	identityverification "kego.com/infrastructure/identity_verification"
	"kego.com/infrastructure/ipresolver"
	"kego.com/infrastructure/logger"
	pushnotification "kego.com/infrastructure/messaging/push_notifications"
	paymentprocessor "kego.com/infrastructure/payment_processor"
)

// Used to start services such as loggers, databases, queues, etc.
func StartServices(){
	logger.InitializeLogger()
	database.SetUpDatabase()
	logger.MetricMonitor.Init()
	logger.RequestMetricMonitor.Init()
	// pubsub.PubSub.Connect()
	fileupload.InitialiseFileUploader()
	pushnotification.InitialisePushNotificationService()
	identityverification.InitialiseIdentityVerifier()
	paymentprocessor.LocalPaymentProcessor.InitialisePaymentProcessor()
	paymentprocessor.InternationalPaymentProcessor.InitialisePaymentProcessor()
	ipresolver.IPResolverInstance.ConnectToDB()
	// sms.InitSMSService()
	application.DBGenesis()
}

// Used to clean up after services that have been shutdown.
func CleanUpServices(){
	datastore.CleanUp()
	// metrics.MetricMonitor.CleanUp()
}