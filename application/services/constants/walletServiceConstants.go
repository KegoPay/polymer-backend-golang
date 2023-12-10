package wallet_constants


type FrozenAccountReason string

var (
	TransactionPinTriesExceeded FrozenAccountReason = "transaction_pin_tries_maxed"
	MaliciousActivitiesDetected FrozenAccountReason = "malicious_activity_detected"
)

type FrozenAccountTime uint

var (
	Min FrozenAccountTime = 60 * 60 * 24 	 // 1 day
	Mid FrozenAccountTime = 60 * 60 * 24 * 3 // 3 day
	Max FrozenAccountTime = 60 * 60 * 24 * 7 // 1 week
)
