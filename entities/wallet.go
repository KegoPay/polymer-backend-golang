package entities

import (
	"time"

	"kego.com/application/utils"
)

type LockedFunds struct {
	Amount               uint64            	`bson:"amount" json:"amount" validate:"required"`
	Reason           	 TransactionIntent  `bson:"reason" json:"reason" validate:"required"`
	LockedFundsID        string    		  	`bson:"lockedFundsID" json:"lockedFundsID" validate:"required"`
	LockedAt 			 time.Time 			`bson:"lockedAt" json:"lockedAt" validate:"required"`
}

type Wallet struct {
	UserID          	string   		 `bson:"userID" json:"userID" validate:"required"`
	BusinessID      	*string   		 `bson:"businessID" json:"businessID"`
	BusinessName      	*string   		 `bson:"businessName" json:"businessName"`
	Frozen          	bool     		 `bson:"frozen" json:"frozen"`
	LedgerBalance 		uint64   		 `bson:"ledgerBalance" json:"-"`
	Balance         	uint64    	 	 `bson:"balance" json:"balance"`
	Currency         	string   		 `bson:"currency" json:"currency" validate:"iso4217"`
	AccountNumber 		*string   		 `bson:"accountNumber" json:"accountNumber"`
	BankName 			*string   		 `bson:"bankName" json:"bankName"`
	LockedFundsLog      []LockedFunds    `bson:"lockedFundsLog" json:"-"`

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
