package startup

import (
	"usepolymer.co/application"
	"usepolymer.co/infrastructure/background"
	"usepolymer.co/infrastructure/biometric"
	cac_service "usepolymer.co/infrastructure/cac"
	"usepolymer.co/infrastructure/database"
	"usepolymer.co/infrastructure/database/connection/datastore"
	fileupload "usepolymer.co/infrastructure/file_upload"
	identityverification "usepolymer.co/infrastructure/identity_verification"
	"usepolymer.co/infrastructure/ipresolver"
	"usepolymer.co/infrastructure/logger"
	pushnotification "usepolymer.co/infrastructure/messaging/push_notifications"
	sms "usepolymer.co/infrastructure/messaging/whatsapp"
	paymentprocessor "usepolymer.co/infrastructure/payment_processor"
	"usepolymer.co/infrastructure/services"
)

// Used to start services such as loggers, databases, queues, etc.
func StartServices() {
	logger.InitializeLogger()
	database.SetUpDatabase()
	// logger.MetricMonitor.Init()
	logger.RequestMetricMonitor.Init()
	// pubsub.PubSub.Connect()
	fileupload.InitialiseFileUploader()
	pushnotification.InitialisePushNotificationService()
	identityverification.InitialiseIdentityVerifier()
	paymentprocessor.LocalPaymentProcessor.InitialisePaymentProcessor()
	paymentprocessor.InternationalPaymentProcessor.InitialisePaymentProcessor()
	ipresolver.IPResolverInstance.ConnectToDB()
	sms.InitSMSService()
	application.DBGenesis()
	biometric.InitialiseBiometricService()
	background.StartScheduler()
	services.InitialiseBackgroundService()
	cac_service.CreateCACServiceInstance()
}

// Used to clean up after services that have been shutdown.
func CleanUpServices() {
	datastore.CleanUp()
	// metrics.MetricMonitor.CleanUp()
}
