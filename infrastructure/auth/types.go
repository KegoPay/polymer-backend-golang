package auth

import "kego.com/entities"

type ClaimsData struct {
    Issuer     string
    UserID     string
    Email      *string
    Phone      *entities.PhoneNumber
    ExpiresAt  int64
    IssuedAt   int64
    DeviceType entities.DeviceType
    DeviceID   string
}
