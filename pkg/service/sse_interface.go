package service

import "github.com/abgeo/mailtm/pkg/dto"

type SSEServiceInterface interface {
	SubscribeMessages(accountID string, handler func(message dto.MessagesItem)) error
}
