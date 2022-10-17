package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/errors"
	"github.com/abgeo/mailtm/pkg/types"
	"github.com/go-resty/resty/v2"
)

type APIService struct {
	client *resty.Client
}

func NewAPIService(version types.Version) *APIService {
	const timeout = 30 * time.Second

	client := resty.New()
	client.SetBaseURL("https://api.mail.tm").
		SetTimeout(timeout).
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", fmt.Sprintf("abgeo-mailtm/%s", version.Number)).
		SetError(&errors.HTTPError{})

	client.JSONMarshal = json.Marshal
	client.JSONUnmarshal = json.Unmarshal

	return &APIService{
		client: client,
	}
}

func (svc *APIService) CreateAccount(data dto.AccountWrite) (account *dto.Account, err error) {
	resp, err := svc.client.R().
		SetBody(data).
		SetResult(&account).
		Post("/accounts")
	if err != nil {
		return account, err
	}

	if resp.StatusCode() != http.StatusCreated {
		return account, errors.NewHTTPError(resp)
	}

	return account, nil
}

func (svc *APIService) GetAccount(id string) (account *dto.Account, err error) {
	resp, err := svc.client.R().
		SetPathParams(types.StrMap{"id": id}).
		SetResult(&account).
		Get("/accounts/{id}")
	if err != nil {
		return account, err
	}

	if resp.StatusCode() != http.StatusOK {
		return account, errors.NewHTTPError(resp)
	}

	return account, nil
}

func (svc *APIService) RemoveAccount(id string) (err error) {
	resp, err := svc.client.R().
		SetPathParams(types.StrMap{"id": id}).
		Delete("/accounts/{id}")
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusNoContent {
		return errors.NewHTTPError(resp)
	}

	return nil
}

func (svc *APIService) GetCurrentAccount() (account *dto.Account, err error) {
	resp, err := svc.client.R().
		SetResult(&account).
		Get("/me")
	if err != nil {
		return account, err
	}

	if resp.StatusCode() != http.StatusOK {
		return account, errors.NewHTTPError(resp)
	}

	return account, nil
}

func (svc *APIService) GetDomains() (domains []dto.Domain, err error) {
	resp, err := svc.client.R().
		SetResult(&domains).
		Get("/domains")
	if err != nil {
		return domains, err
	}

	if resp.StatusCode() != http.StatusOK {
		return domains, errors.NewHTTPError(resp)
	}

	return domains, nil
}

func (svc *APIService) GetDomain(id string) (domain *dto.Domain, err error) {
	resp, err := svc.client.R().
		SetPathParams(types.StrMap{"id": id}).
		SetResult(&domain).
		Get("/domains/{id}")
	if err != nil {
		return domain, err
	}

	if resp.StatusCode() != http.StatusOK {
		return domain, errors.NewHTTPError(resp)
	}

	return domain, nil
}

func (svc *APIService) GetMessages() (messages dto.Messages, err error) {
	resp, err := svc.client.R().
		SetResult(&messages).
		Get("/messages")
	if err != nil {
		return messages, err
	}

	if resp.StatusCode() != http.StatusOK {
		return messages, errors.NewHTTPError(resp)
	}

	return messages, nil
}

func (svc *APIService) GetMessage(id string) (message *dto.Message, err error) {
	resp, err := svc.client.R().
		SetPathParams(types.StrMap{"id": id}).
		SetResult(&message).
		Get("/messages/{id}")
	if err != nil {
		return message, err
	}

	if resp.StatusCode() != http.StatusOK {
		return message, errors.NewHTTPError(resp)
	}

	return message, nil
}

func (svc *APIService) RemoveMessage(id string) (err error) {
	resp, err := svc.client.R().
		SetPathParams(types.StrMap{"id": id}).
		Delete("/messages/{id}")
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusNoContent {
		return errors.NewHTTPError(resp)
	}

	return nil
}

func (svc *APIService) UpdateMessage(id string, data dto.MessageWrite) (err error) {
	resp, err := svc.client.R().
		SetHeader("Content-Type", "application/merge-patch+json").
		SetPathParams(types.StrMap{"id": id}).
		SetBody(data).
		Patch("/messages/{id}")
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.NewHTTPError(resp)
	}

	return nil
}

func (svc *APIService) DownloadMessageAttachment(messageID string, attachmentID string, path string) (err error) {
	resp, err := svc.client.R().
		SetPathParams(types.StrMap{"messageID": messageID, "attachmentID": attachmentID}).
		SetOutput(path).
		Get("/messages/{messageID}/attachment/{attachmentID}")
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.NewHTTPError(resp)
	}

	return nil
}

func (svc *APIService) GetSource(id string) (source *dto.Source, err error) {
	resp, err := svc.client.R().
		SetPathParams(types.StrMap{"id": id}).
		SetResult(&source).
		Get("/sources/{id}")
	if err != nil {
		return source, err
	}

	if resp.StatusCode() != http.StatusOK {
		return source, errors.NewHTTPError(resp)
	}

	return source, nil
}

func (svc *APIService) GetToken(credentials dto.Credentials) (token *dto.Token, err error) {
	resp, err := svc.client.R().
		SetBody(credentials).
		SetResult(&token).
		Post("/token")
	if err != nil {
		return token, err
	}

	if resp.StatusCode() != http.StatusOK {
		return token, errors.NewHTTPError(resp)
	}

	return token, nil
}

func (svc *APIService) SetToken(token *dto.Token) {
	svc.client.SetAuthToken(token.Token)
}
