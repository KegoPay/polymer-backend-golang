package types

type PushNotificationServiceType interface{
	PushOne(deviceID string, header string,  body string)
}