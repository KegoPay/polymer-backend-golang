package entities

import (
	"time"

	"usepolymer.co/application/utils"
	cac_service "usepolymer.co/infrastructure/cac"
)

type ShareHolder struct {
	Name   string `bson:"name" json:"name" validate:"required"`
	ID     string `bson:"id" json:"id"`
	IDType string `bson:"idType" json:"idType"`
	Shares string `bson:"shares" json:"shares"`
}

type Director struct {
	Name   string `bson:"name" json:"name" validate:"required"`
	ID     string `bson:"id" json:"id"`
	IDType string `bson:"idType" json:"idType"`
}

type Business struct {
	Name         string                   `bson:"name" json:"name" validate:"required"`
	UserID       string                   `bson:"userID" json:"userID" validate:"required"`
	WalletID     string                   `bson:"walletID" json:"walletID"`
	Email        string                   `bson:"email" json:"email" validate:"required,email"`
	CACInfo      *cac_service.CACBusiness `bson:"cacInfo" json:"cacInfo"`
	Directors    *[]Director              `bson:"directors" json:"directors"`
	ShareHolders *[]ShareHolder           `bson:"shareholders" json:"shareholders"`

	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

func (business Business) ParseModel() any {
	if business.ID == "" {
		business.CreatedAt = time.Now()
		business.ID = utils.GenerateUUIDString()
	}
	business.UpdatedAt = time.Now()
	return &business
}
