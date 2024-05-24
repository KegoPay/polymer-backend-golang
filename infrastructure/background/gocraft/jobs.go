package gocraft

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gocraft/work"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kego.com/application/repository"
	"kego.com/entities"
	"kego.com/infrastructure/auth"
	"kego.com/infrastructure/cryptography"
	"kego.com/infrastructure/database/repository/cache"
	identityverification "kego.com/infrastructure/identity_verification"
	identity_verification_types "kego.com/infrastructure/identity_verification/types"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/messaging/emails"
	pushnotification "kego.com/infrastructure/messaging/push_notifications"
	"kego.com/infrastructure/services"
)

func SendEmail(job *work.Job) error {
	logger.Info("event processing begining", logger.LoggerOptions{
		Key:  "event_name",
		Data: job.Name,
	})
	success := emails.EmailService.SendEmail(job.ArgString("email"), job.ArgString("subject"), job.ArgString("templateName"), job.Args["opts"])
	if !success {
		logger.Error(errors.New("error sending email in background"), logger.LoggerOptions{
			Key:  "email",
			Data: job.ArgString("email"),
		}, logger.LoggerOptions{
			Key:  "template",
			Data: job.ArgString("templateName"),
		})
	}
	return nil
}

func LockAccount(job *work.Job) error {
	userRepo := repository.UserRepo()
	success, err := userRepo.UpdatePartialByID(job.ArgString("id"), map[string]any{
		"accountLocked": true,
	})
	if err != nil || success == 0 {
		logger.Error(errors.New("error locking user account"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
	}
	return nil
}

func UnlockAccount(job *work.Job) error {
	userRepo := repository.UserRepo()
	success, err := userRepo.UpdatePartialByID(job.ArgString("id"), map[string]any{
		"accountLocked": false,
	})
	if err != nil || success == 0 {
		logger.Error(errors.New("error unlocking user account"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
	}
	return nil
}

func RequestAccountStatement(job *work.Job) error {
	err := services.BackgroundServiceInstance.RequestAccountStatementGeneration(job.ArgString("walletID"), job.ArgString("email"), job.ArgString("start"), job.ArgString("end"))
	if err != nil {
		logger.Error(errors.New("generating account statement failed"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		})
		return err
	}
	return nil
}

func VerifyBusiness(job *work.Job) error {
	businessRepo := repository.BusinessRepo()
	busniess, err := businessRepo.FindByID(job.ArgString("id"), options.FindOne().SetProjection(map[string]any{
		"cacInfo": 1,
		"email":   1,
		"userID":  1,
	}))
	if err != nil {
		logger.Error(errors.New("there was an error fetching business for verification"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "id",
			Data: job.ArgString("id"),
		})
		return err
	}
	if busniess == nil {
		logger.Error(errors.New("business not found"), logger.LoggerOptions{
			Key:  "id",
			Data: job.ArgString("id"),
		})
		return errors.New("buisness not found")
	}
	info, err := identityverification.IdentityVerifier.FetchAdvancedCACDetails(busniess.CACInfo.RCNumber)
	if err != nil {
		logger.Error(errors.New("an error occured while fetching buisiness cac info"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "business",
			Data: busniess,
		})
		return err
	}
	if info == nil {
		logger.Error(errors.New("invalid rc number searched for"), logger.LoggerOptions{
			Key:  "business",
			Data: busniess,
		})
		return errors.New("rc number invalid")
	}
	if os.Getenv("ENV") != "prod" {
		info.Status = "ACTIVE"
	}
	if info.Status == "INACTIVE" {
		logger.Error(errors.New("inactive business profile submitted"), logger.LoggerOptions{
			Key:  "businessID",
			Data: busniess.ID,
		}, logger.LoggerOptions{
			Key:  "info",
			Data: info,
		})
		return errors.New("inactive business profile")
	}
	var match bool
	for _, aff := range info.Affiliates {
		if aff.AffiliateType != "DIRECTOR" {
			continue
		}
		if aff.Email == &busniess.Email {
			match = true
			break
		}
	}
	var user *entities.User
	if !match {
		userRepo := repository.UserRepo()
		user, err = userRepo.FindOneByFilter(map[string]interface{}{
			"email": busniess.Email,
		})
		if err != nil {
			logger.Error(errors.New("there was an error fetching user for business verification"), logger.LoggerOptions{
				Key:  "error",
				Data: err,
			}, logger.LoggerOptions{
				Key:  "id",
				Data: job.ArgString("id"),
			}, logger.LoggerOptions{
				Key:  "info",
				Data: info,
			})
			return err
		}
		if user == nil {
			logger.Error(errors.New("user fetched with business email is not found"), logger.LoggerOptions{
				Key:  "id",
				Data: job.ArgString("id"),
			}, logger.LoggerOptions{
				Key:  "info",
				Data: info,
			}, logger.LoggerOptions{
				Key:  "businessID",
				Data: busniess.ID,
			})
			return errors.New("user fetched with business email is not found")
		}
		for _, aff := range info.Affiliates {
			if aff.AffiliateType != "DIRECTOR" {
				continue
			}
			var plainNIN string
			if user.NIN != "" {
				plainNIN, _ = cryptography.DecryptData(user.NIN, nil)
			}
			plainBVN, _ := cryptography.DecryptData(user.BVN, nil)
			if aff.Email == &user.Email || aff.IDNumber == plainNIN || aff.IDNumber == plainBVN {
				match = true
				break
			}
		}
	}
	fullAddress := fmt.Sprintf("%s, %s, %s", info.City, info.LGA, info.State)
	if match {
		var directors []entities.Director
		var shareholders []entities.ShareHolder
		for _, aff := range info.Affiliates {
			if aff.AffiliateType == "DIRECTOR" {
				directors = append(directors, entities.Director{
					Name:   strings.Replace(aff.Name, ",", "", -1),
					ID:     aff.IDNumber,
					IDType: aff.IDType,
				})
			}
			if aff.AffiliateType == "SHAREHOLDER" {
				shareholders = append(shareholders, entities.ShareHolder{
					ID:     aff.IDNumber,
					IDType: aff.IDType,
					Name:   strings.Replace(aff.Name, ",", "", -1),
					Shares: aff.ShareAllotted,
				})
			}
		}
		updated, err := businessRepo.UpdatePartialByID(busniess.ID, map[string]any{
			"cacInfo.verified":    true,
			"cacInfo.fulladdress": fullAddress,
			"directors":           directors,
			"shareholders":        shareholders,
		})
		if err != nil {
			logger.Error(errors.New("an error occured while updating cac directors and shareholders info"), logger.LoggerOptions{
				Key:  "err",
				Data: err,
			}, logger.LoggerOptions{
				Key:  "info",
				Data: info,
			}, logger.LoggerOptions{
				Key:  "businessID",
				Data: busniess.ID,
			})
			return errors.New("an error occured while updating cac directors and shareholders info")
		}
		if updated != 1 {
			logger.Error(errors.New("could not update business with shareholder and director info"), logger.LoggerOptions{
				Key:  "updated",
				Data: updated,
			}, logger.LoggerOptions{
				Key:  "info",
				Data: info,
			}, logger.LoggerOptions{
				Key:  "businessID",
				Data: busniess.ID,
			})
			return errors.New("could not update business with shareholder and director info")
		}
		pushnotification.PushNotificationService.PushOne(user.PushNotificationToken, "Your business has been verified!🥳", "That was easy, wasn't it? Now you're 1 step closer to unlimited transfer amounts.")
	} else {
		logger.Info("failed to verify business using emails, nin and bvn. Attempting to send emails to directors")
		var directors []entities.Director
		var shareholders []entities.ShareHolder
		for _, aff := range info.Affiliates {
			if aff.AffiliateType == "DIRECTOR" {
				directors = append(directors, entities.Director{
					Name:   strings.Replace(aff.Name, ",", "", -1),
					ID:     aff.IDNumber,
					IDType: aff.IDType,
				})
			}
			if aff.AffiliateType == "SHAREHOLDER" {
				shareholders = append(shareholders, entities.ShareHolder{
					ID:     aff.IDNumber,
					IDType: aff.IDType,
					Name:   strings.Replace(aff.Name, ",", "", -1),
					Shares: aff.ShareAllotted,
				})
			}
		}
		infoByte, _ := json.Marshal(directors)
		success := cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-info-directors", busniess.ID), infoByte, time.Hour*8760) // save for a year
		if !success {
			err = errors.New("an error occured while saving directors information for business for manual review")
			logger.Error(err, logger.LoggerOptions{
				Key:  "info",
				Data: info,
			}, logger.LoggerOptions{
				Key:  "businessID",
				Data: busniess.ID,
			})
			return err
		}
		infoByte, _ = json.Marshal(shareholders)
		success = cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-info-shareholders", busniess.ID), infoByte, time.Hour*8760) // save for a year
		if !success {
			err = errors.New("an error occured while saving shareholders information for business for manual review")
			logger.Error(err, logger.LoggerOptions{
				Key:  "info",
				Data: info,
			}, logger.LoggerOptions{
				Key:  "businessID",
				Data: busniess.ID,
			})
			return err
		}
		token, err := auth.GenerateAuthToken(auth.ClaimsData{
			BusinessID: &busniess.ID,
			UserID:     busniess.UserID,
			ExpiresAt:  time.Now().Local().Add(time.Hour * time.Duration(24*10)).Unix(), //lasts for 10 days
		})
		if err != nil {
			logger.Error(errors.New("an error occured while generating token for business directors to verify manually"), logger.LoggerOptions{
				Key:  "error",
				Data: err,
			}, logger.LoggerOptions{
				Key:  "info",
				Data: info,
			}, logger.LoggerOptions{
				Key:  "businessID",
				Data: busniess.ID,
			})
			return errors.New("an error occured while generating token for business directors to verify manually")
		}
		var wg sync.WaitGroup
		for _, aff := range info.Affiliates {
			if os.Getenv("ENV") != "prod" {
				aff.AffiliateType = "DIRECTOR"
			}
			if aff.AffiliateType != "DIRECTOR" {
				continue
			}
			cache.Cache.CreateEntry(fmt.Sprintf("%s-kyc-info-address", busniess.ID), fullAddress, time.Hour*8760)
			wg.Add(1)
			go func(aff *identity_verification_types.AffiliateProfile) {
				firstName := strings.Split(aff.Name, ",")
				if os.Getenv("ENV") != "prod" {
					aff.Email = &busniess.Email
				}
				success := emails.EmailService.SendEmail(*aff.Email, "Verify your business on Polymer", "verify_business", map[string]string{
					"FIRSTNAME": firstName[0],
					"LINK":      fmt.Sprintf("%s/api/v1/business/verify/manual/%s", os.Getenv("CLIENT_URL"), *token),
				})
				if !success {
					logger.Error(errors.New("error sending email in background"), logger.LoggerOptions{
						Key:  "email",
						Data: job.ArgString("email"),
					}, logger.LoggerOptions{
						Key:  "template",
						Data: job.ArgString("templateName"),
					})
				}
				wg.Done()
			}(&aff)
		}
		wg.Wait()
		pushnotification.PushNotificationService.PushOne(user.PushNotificationToken, "Business verification failed", "A verification email has been sent to the directors. One click from a director verifies the business.")
	}
	return nil
}
