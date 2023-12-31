package types

type PushNotificationServiceType interface{
	PushOne(string, string, string)
}