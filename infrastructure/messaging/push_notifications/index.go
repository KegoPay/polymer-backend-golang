package pushnotification

import (
	firebasepush "usepolymer.co/infrastructure/messaging/push_notifications/firebase"
	"usepolymer.co/infrastructure/messaging/push_notifications/types"
)

var PushNotificationService types.PushNotificationServiceType

func InitialisePushNotificationService() {
	PushNotificationService = (&firebasepush.FireBasePushNotification{}).InitialiseClient()
}
