package gocraft

import (
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
	identityverification "kego.com/infrastructure/identity_verification"
	identity_verification_types "kego.com/infrastructure/identity_verification/types"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/messaging/emails"
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

	if !match {
		userRepo := repository.UserRepo()
		user, err := userRepo.FindOneByFilter(map[string]interface{}{
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
			if aff.Email == &user.Email || aff.IDNumber == user.NIN || aff.IDNumber == user.BVN {
				match = true
				break
			}
		}
	}
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
					IDType: aff.AffiliateType,
					Name:   strings.Replace(aff.Name, ",", "", -1),
					Shares: aff.ShareAllotted,
				})
			}
		}
		updated, err := businessRepo.UpdatePartialByID(busniess.ID, map[string]any{
			"cacInfo.verified":    true,
			"cacInfo.fulladdress": fmt.Sprintf("%s, %s, %s", info.City, info.LGA, info.State),
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
	} else {
		logger.Info("failed to verify business using emails, nin and bvn. Attempting to send emails to directors")
		token, err := auth.GenerateAuthToken(auth.ClaimsData{
			BusinessID: &busniess.ID,
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
		wg.Add(1)
		for _, aff := range info.Affiliates {
			if os.Getenv("ENV") != "prod" {
				aff.AffiliateType = "DIRECTOR"
			}
			if aff.AffiliateType != "DIRECTOR" {
				continue
			}
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

	}
	return nil
}
