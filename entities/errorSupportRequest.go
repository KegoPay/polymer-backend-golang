package entities

import (
	"time"

	"usepolymer.co/application/utils"
)

type ErrorSupportRequest struct {
	UserID   string `bson:"userID" json:"userID" validate:"required"`
	Message  string `bson:"message" json:"message" validate:"required"`
	Email    string `bson:"email" json:"email" validate:"required"`
	Resolved bool   `bson:"resolved" json:"resolved"`

	ID        string    `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

func (sr ErrorSupportRequest) ParseModel() any {
	now := time.Now()
	if sr.ID == "" {
		sr.CreatedAt = now
		sr.ID = utils.GenerateUUIDString()
	}
	sr.UpdatedAt = now
	return &sr
}
