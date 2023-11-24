package authusecases

import (
	"time"

	apperrors "kego.com/application/appErrors"
	"kego.com/application/repository"
	"kego.com/entities"
	"kego.com/infrastructure/auth"
	"kego.com/infrastructure/cryptography"
	"kego.com/infrastructure/database/repository/cache"
)

func LoginAccount(ctx any, email *string, phone *string, password *string) (*entities.User, *string) {
	userRepo := repository.UserRepo()
	var account *entities.User
	var err error
	if email != nil {
		account, err = userRepo.FindOneByFilter(map[string]interface{}{
			"email": email,
		})
	}
	if phone != nil {
		account, err = userRepo.FindOneByFilter(map[string]interface{}{
			"phone.localNumber": phone,
		})
	}
	if err != nil {
		apperrors.FatalServerError(ctx)
		return nil, nil
	}
	if account == nil {
		apperrors.NotFoundError(ctx, "this account does not exist")
		return nil, nil
	}
	if  phone != nil && !account.PhoneVerified {
		apperrors.ClientError(ctx, "verify your phone number to use it to login", nil)
		return nil, nil
	} 
	passwordMatch := cryptography.CryptoHahser.VerifyData(account.Password, *password)
	if !passwordMatch {
		apperrors.AuthenticationError(ctx, "wrong password")
		return nil, nil
	}
	token, err := auth.GenerateAuthToken(auth.ClaimsData{
		Email:     &account.Email,
		Phone:     &account.Phone,
		UserID:    account.ID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(10)).Unix(), //lasts for 10 mins
		DeviceType: account.UserAgent,
		DeviceID:   account.DeviceID,
	})
	cache.Cache.CreateEntry(account.ID, *token, time.Minute * time.Duration(10)) // cache authentication token for 10 mins
	return account, token
}