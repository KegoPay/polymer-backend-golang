package entities

import (
	"time"

	"kego.com/application/utils"
	cac_service "kego.com/infrastructure/cac"
)

type Business struct {
	Name      	string    				   `bson:"name" json:"name" validate:"required"`
	UserID    	string      			   `bson:"userID" json:"userID" validate:"required"`
	WalletID  	string   				   `bson:"walletID" json:"walletID"`
	Email	  	string     				   `bson:"email" json:"email" validate:"required,email"`
	CACInfo	  	*cac_service.CACBusiness    `bson:"cacInfo" json:"cacInfo"`

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
