package emails

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"usepolymer.co/infrastructure/logger"
)

var dir, _ = os.Getwd()

type SendGridService struct {
}

func (s *SendGridService) SendEmail(toEmail string, subject string, templateName string, opts interface{}) bool {
	from := mail.NewEmail("Kego", os.Getenv("KEGO_EMAIL"))
	to := mail.NewEmail(toEmail, toEmail)
	buffer := s.loadTemplates(templateName, opts)
	message := mail.NewSingleEmail(from, subject, to, "", buffer.String())
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		logger.Error(err, logger.LoggerOptions{
			Key:  "to",
			Data: toEmail,
		}, logger.LoggerOptions{
			Key:  "templateName",
			Data: templateName,
		})
		return false
	} else {
		if response.StatusCode == 202 {
			logger.Info(fmt.Sprintf("email sent to %s successfully", toEmail), logger.LoggerOptions{
				Key:  "response",
				Data: response,
			})
			return true
		} else {
			logger.Info(fmt.Sprintf("failed to send email to %s", toEmail), logger.LoggerOptions{
				Key:  "response",
				Data: response,
			})
			return false
		}
	}
}

func (s *SendGridService) loadTemplates(templateName string, opts interface{}) bytes.Buffer {
	var buffer bytes.Buffer
	template.Must(template.ParseFiles(fmt.Sprintf(filepath.Join(dir, "/infrastructure/messaging/emails/templates/%s.html"), templateName))).Execute(&buffer, opts)
	return buffer
}
