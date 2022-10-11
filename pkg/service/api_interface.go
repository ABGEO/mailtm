//nolint:interfacebloat
package service

import (
	"github.com/abgeo/mailtm/pkg/dto"
)

type APIServiceInterface interface {
	CreateAccount(data dto.AccountWrite) (account *dto.Account, err error)
	GetAccount(id string) (account *dto.Account, err error)
	RemoveAccount(id string) (err error)
	GetCurrentAccount() (account *dto.Account, err error)
	GetDomains() (domains []dto.Domain, err error)
	GetDomain(id string) (domain *dto.Domain, err error)
	GetMessages() (messages dto.Messages, err error)
	GetMessage(id string) (message *dto.Message, err error)
	RemoveMessage(id string) (err error)
	UpdateMessage(id string, data dto.MessageWrite) (err error)
	DownloadMessageAttachment(messageID string, attachmentID string, path string) (err error)
	GetSource(id string) (source *dto.Source, err error)
	GetToken(credentials dto.Credentials) (token *dto.Token, err error)
	SetToken(token *dto.Token)
}
