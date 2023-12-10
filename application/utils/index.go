package utils

import (
	"github.com/google/uuid"
)

func GenerateUUIDString() string {
	return uuid.NewString()
}

func ParseAmountToSmallerDenomination(amount uint64) uint64 {
	return amount * 100
}

func ParseAmountToHigherDenomination(amount uint64) uint64 {
	return amount / 100
}