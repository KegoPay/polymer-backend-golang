package emails

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/resend/resend-go/v2"
	"usepolymer.co/infrastructure/logger"
)

var rsdir, _ = os.Getwd()

type ResendService struct {
}

func (rs *ResendService) SendEmail(toEmail string, subject string, templateName string, opts interface{}) bool {
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	html := rs.loadTemplates(templateName, opts)

	params := &resend.SendEmailRequest{
		From:    os.Getenv("POLYMER_DEFAULT_EMAIL"),
		To:      []string{toEmail},
		Subject: subject,
		Html:    *html,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		logger.Error(errors.New("an error occured while trying to send email using resend service"), logger.LoggerOptions{
			Key:  "error",
			Data: err,
		}, logger.LoggerOptions{
			Key:  "toEmail",
			Data: toEmail,
		}, logger.LoggerOptions{
			Key:  "templateName",
			Data: templateName,
		})
		return false
	}
	logger.Info(fmt.Sprintf("successfully sent email to %s", toEmail), logger.LoggerOptions{
		Key:  "templateName",
		Data: templateName,
	}, logger.LoggerOptions{
		Key:  "service",
		Data: "resend",
	})
	return true
}

func (rs *ResendService) loadTemplates(templateName string, opts interface{}) *string {
	var buffer bytes.Buffer
	template.Must(template.ParseFiles(fmt.Sprintf(filepath.Join(rsdir, "/infrastructure/messaging/emails/templates/%s.html"), templateName))).Execute(&buffer, opts)
	templateString := buffer.String()
	return &templateString
}
