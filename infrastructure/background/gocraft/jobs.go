package gocraft

import (
	"errors"

	"github.com/gocraft/work"
	"kego.com/infrastructure/logger"
	"kego.com/infrastructure/messaging/emails"
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