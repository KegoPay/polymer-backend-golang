package pushnotification

import (
	firebase_push_notification "kego.com/infrastructure/push_notification/firebase"
	"kego.com/infrastructure/push_notification/types"
)

var PushNotificationService types.PushNotificationServiceType = (&firebase_push_notification.FireBasePushNotification{}).InitialiseClient()