package mail

import (
	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/nuricanozturk01/kariyerklubu-lib/mail/credentials"
	"github.com/nuricanozturk01/kariyerklubu-lib/mail/form"
)

type Mail struct {
	MailjetClient *mailjet.Client
}

func NewMailClient(credentials *credentials.MailCredentials) *Mail {
	return &Mail{
		MailjetClient: getMailClient(credentials),
	}
}

func getMailClient(credentials *credentials.MailCredentials) *mailjet.Client {
	return mailjet.NewMailjetClient(credentials.ApiKey, credentials.SecretKey)
}

func (m *Mail) SendEmailTemplate(emailForm *form.EmailTemplateForm) (string, error) {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: form.KariyerKlubuEmail,
				Name:  form.KariyerKlubuName,
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
				Email: form.KariyerKlubuEmail,
				Name:  form.KariyerKlubuName,
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
