package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/nuricanozturk01/kariyerklubu-lib/slack/credentials"
	"io"
	"net/http"
	"time"
)

const (
	SlackWarn    string = "#ffcc00"
	SlackSuccess string = "#36a64f"
	SlackError   string = "#ff0000"
	SlackInfo    string = "#00008b"
)

type Slack struct {
	slackCredentials *credentials.SlackCredentials
}

func NewSlack(credentials *credentials.SlackCredentials) *Slack {
	return &Slack{
		slackCredentials: credentials,
	}
}

func (s *Slack) sendNotification(hasMarkdown bool, message, messageType string) error {
	currentTime := time.Now().Format(time.DateTime)
	message = fmt.Sprintf("%s - %s ", currentTime, message)
	payload := map[string]any{
		"attachments": []map[string]any{
			{
				"color": messageType,
				"text":  message,
			},
		},
	}

	if hasMarkdown {
		message = fmt.Sprintf("*%s* - %s", currentTime, message)
		payload["attachments"].([]map[string]any)[0]["mrkdwn_in"] = []string{"text"}
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error creating Slack payload: %v", err)
	}

	resp, err := http.Post(s.slackCredentials.WebHookUrl, "application/json", bytes.NewBuffer(payloadBytes))

	if err != nil {
		return fmt.Errorf("error sending Slack message: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("Error closing response body: ", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response from Slack: %v", resp.Status)
	}

	return nil
}

func (s *Slack) SendSlackMessageMarkdown(message, messageType string) {
	if s.slackCredentials.Enable {
		if err := s.sendNotification(true, message, messageType); err != nil {
			fmt.Println("Failed to send slack message")
		}
	}
}

func (s *Slack) SendSlackMessage(message, messageType string) {
	if s.slackCredentials.Enable {
		if err := s.sendNotification(false, message, messageType); err != nil {
			fmt.Println("Failed to send slack message")
		}
	}
}
