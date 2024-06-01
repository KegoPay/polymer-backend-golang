package authusecases

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	apperrors "usepolymer.co/application/appErrors"
	"usepolymer.co/application/constants"
	"usepolymer.co/application/repository"
	"usepolymer.co/application/utils"
	"usepolymer.co/entities"
	"usepolymer.co/infrastructure/auth"
	"usepolymer.co/infrastructure/background"
	"usepolymer.co/infrastructure/cryptography"
	"usepolymer.co/infrastructure/database/repository/cache"
	geospatialcalculation "usepolymer.co/infrastructure/geospatial_calculation"
	"usepolymer.co/infrastructure/logger"
)

func LoginAccount(ctx any, email *string, phone *string, password *string, appVersion string, userAgent string, deviceID string, pushNotificationToken string, longitude float64, latitude float64, device_id *string) (*entities.User, *entities.Wallet, *string) {
	currentTries := cache.Cache.FindOne(fmt.Sprintf("%s-password-tries", *email))
	if currentTries == nil {
		currentTries = utils.GetStringPointer("0")
	}
	currentTriesInt, err := strconv.Atoi(*currentTries)
	if err != nil {
		logger.Error(errors.New("error parsing users password current tries"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "key",
			Data: fmt.Sprintf("%s-password-tries", *email),
		}, logger.LoggerOptions{
			Key:  "data",
			Data: currentTries,
		})
		apperrors.FatalServerError(ctx, err, device_id)
		return nil, nil, nil
	}
	if currentTriesInt == constants.MAX_PASSWORD_TRIES {
		err = errors.New("you have exceeded the number of tries for your password and your account has been temporarily locked for 5 days")
		apperrors.AuthenticationError(ctx, err.Error(), device_id)
		return nil, nil, nil
	}
	userRepo := repository.UserRepo()
	var account *entities.User
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
		apperrors.FatalServerError(ctx, err, device_id)
		return nil, nil, nil
	}
	if account == nil {
		apperrors.NotFoundError(ctx, "this account does not exist", device_id)
		return nil, nil, nil
	}
	if !account.EmailVerified {
		apperrors.ClientError(ctx, "verify your email to use it to login", nil, &constants.EMAIL_UNVERIFIED, device_id)
		return nil, nil, nil
	}
	if account.Deactivated {
		apperrors.ClientError(ctx, "this account has been deactivated", nil, nil, device_id)
		return nil, nil, nil
	}
	passwordMatch := cryptography.CryptoHahser.VerifyData(account.Password, *password)
	if !passwordMatch {
		currentTriesInt = currentTriesInt + 1
		cache.Cache.CreateEntry(fmt.Sprintf("%s-password-tries", *email), fmt.Sprintf("%d", currentTriesInt), time.Hour*24*5)
		msg := fmt.Sprintf("wrong password. your account will be deactivated after %d wrong attempts", constants.MAX_PASSWORD_TRIES-currentTriesInt)
		if currentTriesInt == 0 {
			msg = "you have exceeded maximum password tries and your account has been locked"
		}
		apperrors.AuthenticationError(ctx, msg, device_id)
		return nil, nil, nil
	}
	cache.Cache.CreateEntry(fmt.Sprintf("%s-password-tries", *email), 0, 0)
	distance := geospatialcalculation.GeospatialCalculator.CalculateDistanceInMiles(account.Longitude, longitude, account.Latitude, latitude)
	if *distance > 2 {
		logger.Warning("this triggers a wallet lock")
		background.Scheduler.Emit("lock_account", map[string]any{
			"id": account.ID,
		})
	}
	token, err := auth.GenerateAuthToken(auth.ClaimsData{
		Email:                 &account.Email,
		Phone:                 account.Phone,
		UserID:                account.ID,
		IssuedAt:              time.Now().Unix(),
		ExpiresAt:             time.Now().Local().Add(time.Minute * time.Duration(10)).Unix(), //lasts for 10 mins
		UserAgent:             account.UserAgent,
		FirstName:             account.FirstName,
		LastName:              account.LastName,
		DeviceID:              account.DeviceID,
		AppVersion:            account.AppVersion,
		PushNotificationToken: account.PushNotificationToken,
	})
	if err != nil {
		apperrors.UnknownError(ctx, err, device_id)
		return nil, nil, nil
	}
	var updateAccountPayload = map[string]any{}
	if account.UserAgent != userAgent {
		updateAccountPayload["userAgent"] = userAgent
		account.UserAgent = userAgent
	}
	if appVersion != account.AppVersion {
		updateAccountPayload["appVersion"] = appVersion
		account.AppVersion = appVersion
	}
	updateAccountPayload["deviceID"] = deviceID
	updateAccountPayload["pushNotificationToken"] = pushNotificationToken
	account.DeviceID = deviceID
	userRepo.UpdatePartialByID(account.ID, updateAccountPayload)
	hashedToken, err := cryptography.CryptoHahser.HashString(*token)
	if err != nil {
		apperrors.FatalServerError(ctx, err, device_id)
		return nil, nil, nil
	}
	cache.Cache.CreateEntry(account.ID, hashedToken, time.Minute*time.Duration(10)) // cache authentication token for 10 mins
	walletRepo := repository.WalletRepo()
	wallet, err := walletRepo.FindByID(account.WalletID)
	if err != nil {
		apperrors.FatalServerError(ctx, fmt.Errorf("failed to find user wallet by id on login. walletID %s", account.WalletID), device_id)
		return nil, nil, nil
	}
	return account, wallet, token
}
