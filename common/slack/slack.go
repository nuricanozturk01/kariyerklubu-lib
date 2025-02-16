package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/nuricanozturk01/kariyerklubu-lib/common/config"

	"io"
	"net/http"
	"time"
)

type Slack struct {
	configuration *config.Config
}

func NewSlack(configuration *config.Config) *Slack {
	return &Slack{
		configuration: configuration,
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

	resp, err := http.Post(s.configuration.SlackWebhookURL, "application/json", bytes.NewBuffer(payloadBytes))

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
	if s.configuration.EnableSlackNotifications {
		if err := s.sendNotification(true, message, messageType); err != nil {
			fmt.Println("Failed to send slack message")
		}
	}
}

func (s *Slack) SendSlackMessage(message, messageType string) {
	if s.configuration.EnableSlackNotifications {
		if err := s.sendNotification(false, message, messageType); err != nil {
			fmt.Println("Failed to send slack message")
		}
	}
}
