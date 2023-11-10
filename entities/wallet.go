package entities

import (
	"time"

	"kego.com/application/utils"
)

type Wallet struct {
	UserID          	string    `bson:"userID" json:"userID" validate:"required"`
	BusinessID      	string    `bson:"businessID" json:"businessID" validate:"required"`
	Frozen          	bool      `bson:"frozen" json:"frozen"`
	AvailableBalance 	uint      `bson:"availableBalance" json:"availableBalance"`
	Balance         	uint      `bson:"balance" json:"balance"`

	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

func (wallet Wallet) ParseModel() any {
	if wallet.ID == "" {
		wallet.CreatedAt = time.Now()
		wallet.ID = utils.GenerateUUIDString()
	}
	wallet.UpdatedAt = time.Now()
	return &wallet
}
