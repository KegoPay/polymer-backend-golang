package pushnotification

import (
	firebasepush "kego.com/infrastructure/messaging/push_notifications/firebase"
	"kego.com/infrastructure/messaging/push_notifications/types"
)

var PushNotificationService types.PushNotificationServiceType 

func InitialisePushNotificationService() {
	PushNotificationService = (&firebasepush.FireBasePushNotification{}).InitialiseClient()
}