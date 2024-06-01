package entities

import (
	"time"

	wallet_constants "usepolymer.co/application/services/constants"
	"usepolymer.co/application/utils"
)

type FrozenWalletLog struct {
	WalletID string                               `bson:"walletID" json:"walletID"`
	UserID   string                               `bson:"userID" json:"userID"`
	Reason   wallet_constants.FrozenAccountReason `bson:"reason" json:"reason"`
	Time     wallet_constants.FrozenAccountTime   `bson:"time" json:"time"`
	Unfrozen bool                                 `bson:"unfrozen" json:"unfrozen"`

	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

func (fwLog FrozenWalletLog) ParseModel() any {
	if fwLog.ID == "" {
		fwLog.CreatedAt = time.Now()
		fwLog.ID = utils.GenerateUUIDString()
	}
	fwLog.UpdatedAt = time.Now()
	return &fwLog
}
