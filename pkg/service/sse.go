package service

import (
	"encoding/json"
	"fmt"

	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/types"
	"github.com/r3labs/sse/v2"
)

type SSEService struct {
	client *sse.Client
}

func NewSSEService(version types.Version, auth configs.AuthConfig) *SSEService {
	client := sse.NewClient("https://mercure.mail.tm/.well-known/mercure", func(c *sse.Client) {
		c.Headers = map[string]string{
			"Accept":        "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", auth.Token),
			"User-Agent":    fmt.Sprintf("abgeo-mailtm/%s", version.Number),
		}
	})

	return &SSEService{
		client: client,
	}
}

func (svc *SSEService) SubscribeMessages(accountID string, handler func(message dto.MessagesItem)) error {
	// Hack to override URL query parameter.
	// SSE uses "stream" as a key, but we need "topic"
	svc.client.URL = fmt.Sprintf("%s?topic=/accounts/%s", svc.client.URL, accountID)

	return svc.client.Subscribe("", func(msg *sse.Event) {
		var message dto.MessagesItem

		if err := json.Unmarshal(msg.Data, &message); err == nil {
			if message.AccountID != "" {
				handler(message)
			}
		}
	})
}
