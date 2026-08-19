package auth

import (
	"bytes"
	"html/template"
	"log"
	"os"

	"github.com/Sushank-ghimire/url-shortener/internal/config"
	"github.com/resend/resend-go/v3"
)

type MailParams struct {
	From    string
	To      []string
	HTML    string
	Subject string
}

func SendMail(resendParams MailParams, config *config.Config) (sentId string, error error) {
	client := resend.NewClient(config.ResendKey)

	params := &resend.SendEmailRequest{
		From:    config.FromEmail,
		To:      resendParams.To,
		Html:    resendParams.HTML,
		Subject: resendParams.Subject,
	}
	sent, err := client.Emails.Send(params)

	if err != nil {
		log.Printf("error sending mail: %v", err)
		return "", err
	}

	return sent.Id, nil
}

func LoadTemplate(path string, data any) (string, error) {
	templ, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = templ.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func SendVerificationEmail(
	to string,
	name string,
	verificationURL string,
	config *config.Config,
) (string, error) {
	html, err := LoadTemplate(
		"templates/mail/verification.html",
		struct {
			Name            string
			VerificationURL string
		}{
			Name:            name,
			VerificationURL: verificationURL,
		},
	)

	if err != nil {
		return "", err
	}
	return SendMail(
		MailParams{
			From:    os.Getenv(""),
			To:      []string{to},
			HTML:    html,
			Subject: "Verify your email",
		},
		config,
	)
}
