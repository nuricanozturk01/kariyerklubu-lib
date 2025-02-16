package mail

import (
	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/nuricanozturk01/kariyerklubu-lib/common/config"
	"github.com/nuricanozturk01/kariyerklubu-lib/common/mail/dto"
)

type Mail struct {
	configuration *config.Config
	MailjetClient *mailjet.Client
}

func NewMailClient(configuration *config.Config) *Mail {
	return &Mail{
		configuration: configuration,
		MailjetClient: getMailClient(configuration),
	}
}

func getMailClient(configuration *config.Config) *mailjet.Client {
	return mailjet.NewMailjetClient(configuration.MailjetAPIKey, configuration.MailjetSecretKey)
}

func (m *Mail) SendEmailTemplate(emailForm *dto.EmailTemplateForm) (string, error) {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: dto.KariyerKlubuEmail,
				Name:  dto.KariyerKlubuName,
			},
			To: &mailjet.RecipientsV31{
				{
					Email: emailForm.To,
					Name:  emailForm.Name,
				},
			},
			TemplateLanguage: true,
			TemplateID:       emailForm.TemplateID,
			Variables:        emailForm.Variables,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}

	res, err := m.MailjetClient.SendMailV31(&messages)
	if err != nil {
		return "failed", err
	}

	return res.ResultsV31[0].Status, nil
}

func (m *Mail) SendEmailStr(subject, message, to, name string) (string, error) {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: dto.KariyerKlubuEmail,
				Name:  dto.KariyerKlubuName,
			},
			To: &mailjet.RecipientsV31{
				{
					Email: to,
					Name:  name,
				},
			},
			TemplateLanguage: true,
			Subject:          subject,
			TextPart:         message,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}

	res, err := m.MailjetClient.SendMailV31(&messages)
	if err != nil {
		return "failed", err
	}

	return res.ResultsV31[0].Status, nil
}
