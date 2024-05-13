package gocraft

import (
	"errors"

	"github.com/gocraft/work"
	"kego.com/application/repository"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/messaging/emails"
	"kego.com/infrastructure/services"
)

func SendEmail(job *work.Job) error {
	logger.Info("event processing begining", logger.LoggerOptions{
		Key: "event_name",
		Data: job.Name,
	})
	success := emails.EmailService.SendEmail(job.ArgString("email"), job.ArgString("subject"), job.ArgString("templateName"), job.Args["opts"])
	if !success {
		logger.Error(errors.New("error sending email in background"), logger.LoggerOptions{
			Key: "email",
			Data: job.ArgString("email"),
		}, logger.LoggerOptions{
			Key: "template",
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
			Key: "error",
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
			Key: "error",
			Data: err,
		})
	}
	return nil
}

func RequestAccountStatement(job *work.Job) error {
	err := services.BackgroundServiceInstance.RequestAccountStatementGeneration(job.ArgString("walletID"), job.ArgString("email"), job.ArgString("start"), job.ArgString("end"))
	if err != nil {
		logger.Error(errors.New("generating account statement failed"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return err
	}
	return nil
}