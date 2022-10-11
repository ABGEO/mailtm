package service

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/errors"
	"github.com/abgeo/mailtm/pkg/types"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type APIServiceSuite struct {
	suite.Suite

	Service *APIService
}

func (suite *APIServiceSuite) SetupSuite() {
	version := types.Version{
		Number: "test",
	}
	suite.Service = NewAPIService(version)

	httpmock.ActivateNonDefault(suite.Service.client.GetClient())
}

func (suite *APIServiceSuite) TearDownTest() {
	httpmock.Reset()

	err := os.RemoveAll("/tmp/attachments")
	if err != nil {
		log.Fatal(err)
	}
}

func (suite *APIServiceSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func TestAPIServiceSuite(t *testing.T) {
	suite.Run(t, new(APIServiceSuite))
}

func (suite *APIServiceSuite) TestUserAgent() {
	assert.Equal(suite.T(), "abgeo-mailtm/test", suite.Service.client.Header.Get("User-Agent"))
}

func (suite *APIServiceSuite) TestCreateAccount_Success() {
	fixture := map[string]interface{}{
		"id":         "acc001",
		"address":    "foo@bar.com",
		"quota":      100,
		"used":       0,
		"isDisabled": false,
		"isDeleted":  false,
		"createdAt":  "2022-10-06T08:59:36.792Z",
		"updatedAt":  "2022-10-06T08:59:36.792Z",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusCreated, fixture)
	httpmock.RegisterResponder(http.MethodPost, "/accounts", responder)

	body := dto.AccountWrite{
		Address:  "foo@bar.com",
		Password: "Pa$$w0rd",
	}
	response, err := suite.Service.CreateAccount(body)

	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), &dto.Account{}, response)
	assert.Equal(suite.T(), "acc001", response.ID)
	assert.Equal(suite.T(), false, response.IsDeleted)
	assert.Equal(suite.T(), 100, response.Quota)
}

func (suite *APIServiceSuite) TestCreateAccount_InvalidInput() {
	fixture := map[string]any{
		"type":   "https://tools.ietf.org/html/rfc2616#section-10",
		"title":  "An error occurred",
		"detail": "address: This value is not a valid email address.",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusBadRequest, fixture)
	httpmock.RegisterResponder(http.MethodPost, "/accounts", responder)

	body := dto.AccountWrite{
		Address:  "",
		Password: "",
	}
	response, err := suite.Service.CreateAccount(body)

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.IsType(suite.T(), &errors.HTTPError{}, err)
	assert.Equal(suite.T(), "address: This value is not a valid email address. [400]", err.Error())
}

func (suite *APIServiceSuite) TestCreateAccount_InvalidResponse() {
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, "invalid data")
	httpmock.RegisterResponder(http.MethodPost, "/accounts", responder)

	_, err := suite.Service.CreateAccount(dto.AccountWrite{})

	assert.IsType(suite.T(), &json.UnmarshalTypeError{}, err)
}

func (suite *APIServiceSuite) TestGetAccount_Success() {
	fixture := map[string]any{
		"id":         "acc001",
		"address":    "foo@bar.com",
		"quota":      100,
		"used":       0,
		"isDisabled": false,
		"isDeleted":  false,
		"createdAt":  "2022-10-06T08:59:36.792Z",
		"updatedAt":  "2022-10-06T08:59:36.792Z",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusOK, fixture)
	httpmock.RegisterResponder(http.MethodGet, `=~/accounts/\w+\z`, responder)

	response, err := suite.Service.GetAccount("acc001")

	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), &dto.Account{}, response)
	assert.Equal(suite.T(), "acc001", response.ID)
}

func (suite *APIServiceSuite) TestGetAccount_NotFound() {
	fixture := map[string]any{
		"type":   "https://tools.ietf.org/html/rfc2616#section-10",
		"title":  "An error occurred",
		"detail": "Not Found",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusNotFound, fixture)
	httpmock.RegisterResponder(http.MethodGet, `=~/accounts/\w+\z`, responder)

	response, err := suite.Service.GetAccount("acc001")

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.IsType(suite.T(), &errors.HTTPError{}, err)
	assert.Equal(suite.T(), "Not Found [404]", err.Error())
}

func (suite *APIServiceSuite) TestGetAccount_InvalidResponse() {
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, "invalid data")
	httpmock.RegisterResponder(http.MethodGet, `=~/accounts/\w+\z`, responder)

	_, err := suite.Service.GetAccount("acc001")

	assert.IsType(suite.T(), &json.UnmarshalTypeError{}, err)
}

func (suite *APIServiceSuite) TestRemoveAccount_Success() {
	fixture := map[string]any{}
	responder, _ := httpmock.NewJsonResponder(http.StatusNoContent, fixture)
	httpmock.RegisterResponder(http.MethodDelete, `=~/accounts/\w+\z`, responder)

	err := suite.Service.RemoveAccount("acc001")

	assert.Nil(suite.T(), err)
}

func (suite *APIServiceSuite) TestRemoveAccount_InternalServerError() {
	fixture := map[string]any{
		"type":   "https://tools.ietf.org/html/rfc2616#section-10",
		"title":  "An error occurred",
		"detail": "Internal Server Error",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, fixture)
	httpmock.RegisterResponder(http.MethodDelete, `=~/accounts/\w+\z`, responder)

	err := suite.Service.RemoveAccount("acc001")

	assert.NotNil(suite.T(), err)
	assert.IsType(suite.T(), &errors.HTTPError{}, err)
	assert.Equal(suite.T(), "Internal Server Error [500]", err.Error())
}

func (suite *APIServiceSuite) TestRemoveAccount_InvalidResponse() {
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, "invalid data")
	httpmock.RegisterResponder(http.MethodDelete, `=~/accounts/\w+\z`, responder)

	err := suite.Service.RemoveAccount("acc001")

	assert.IsType(suite.T(), &json.UnmarshalTypeError{}, err)
}

func (suite *APIServiceSuite) TestCurrentAccount_Success() {
	fixture := map[string]any{
		"id":         "acc001",
		"address":    "foo@bar.com",
		"quota":      100,
		"used":       0,
		"isDisabled": false,
		"isDeleted":  false,
		"createdAt":  "2022-10-06T08:59:36.792Z",
		"updatedAt":  "2022-10-06T08:59:36.792Z",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusOK, fixture)
	httpmock.RegisterResponder(http.MethodGet, "/me", responder)

	response, err := suite.Service.GetCurrentAccount()

	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), &dto.Account{}, response)
	assert.Equal(suite.T(), "acc001", response.ID)
}

func (suite *APIServiceSuite) TestCurrentAccount_Unauthorized() {
	fixture := map[string]any{
		"code":    401,
		"message": "JWT Token not found",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusUnauthorized, fixture)
	httpmock.RegisterResponder(http.MethodGet, "/me", responder)

	response, err := suite.Service.GetCurrentAccount()

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.IsType(suite.T(), &errors.HTTPError{}, err)
	assert.Equal(suite.T(), "JWT Token not found [401]", err.Error())
}

func (suite *APIServiceSuite) TestCurrentAccount_InvalidResponse() {
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, "invalid data")
	httpmock.RegisterResponder(http.MethodGet, "/me", responder)

	_, err := suite.Service.GetCurrentAccount()

	assert.IsType(suite.T(), &json.UnmarshalTypeError{}, err)
}

func (suite *APIServiceSuite) TestGetDomains_Success() {
	fixture := []map[string]any{
		{
			"id":        "dom001",
			"domain":    "foo.bar",
			"isActive":  true,
			"createdAt": "2022-10-06T12:57:27.262Z",
			"updatedAt": "2022-10-06T12:57:27.262Z",
		},
		{
			"id":        "dom002",
			"domain":    "foo.bar",
			"isActive":  false,
			"createdAt": "2022-10-06T12:57:27.262Z",
			"updatedAt": "2022-10-06T12:57:27.262Z",
		},
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusOK, fixture)
	httpmock.RegisterResponder(http.MethodGet, "/domains", responder)

	response, err := suite.Service.GetDomains()

	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), []dto.Domain{}, response)
	assert.Equal(suite.T(), 2, len(response))
	assert.Equal(suite.T(), "dom001", response[0].ID)
	assert.False(suite.T(), response[1].IsActive)
}

func (suite *APIServiceSuite) TestGetDomains_InternalServerError() {
	fixture := map[string]any{
		"code":    500,
		"message": "Internal Server Error",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, fixture)
	httpmock.RegisterResponder(http.MethodGet, "/domains", responder)

	response, err := suite.Service.GetDomains()

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.IsType(suite.T(), &errors.HTTPError{}, err)
	assert.Equal(suite.T(), "Internal Server Error [500]", err.Error())
}

func (suite *APIServiceSuite) TestGetDomains_InvalidResponse() {
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, "invalid data")
	httpmock.RegisterResponder(http.MethodGet, "/domains", responder)

	_, err := suite.Service.GetDomains()

	assert.IsType(suite.T(), &json.UnmarshalTypeError{}, err)
}

func (suite *APIServiceSuite) TestGetDomain_Success() {
	fixture := map[string]any{
		"id":        "dom001",
		"domain":    "foo.bar",
		"isActive":  true,
		"createdAt": "2022-10-06T12:57:27.262Z",
		"updatedAt": "2022-10-06T12:57:27.262Z",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusOK, fixture)
	httpmock.RegisterResponder(http.MethodGet, `=~/domains/\w+\z`, responder)

	response, err := suite.Service.GetDomain("dom001")

	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), &dto.Domain{}, response)
	assert.Equal(suite.T(), "dom001", response.ID)
}

func (suite *APIServiceSuite) TestGetDomain_NotFound() {
	fixture := map[string]any{
		"type":   "https://tools.ietf.org/html/rfc2616#section-10",
		"title":  "An error occurred",
		"detail": "Not Found",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusNotFound, fixture)
	httpmock.RegisterResponder(http.MethodGet, `=~/domains/\w+\z`, responder)

	response, err := suite.Service.GetDomain("dom001")

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.IsType(suite.T(), &errors.HTTPError{}, err)
	assert.Equal(suite.T(), "Not Found [404]", err.Error())
}

func (suite *APIServiceSuite) TestGetDomain_InvalidResponse() {
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, "invalid data")
	httpmock.RegisterResponder(http.MethodGet, `=~/domains/\w+\z`, responder)

	_, err := suite.Service.GetDomain("dom001")

	assert.IsType(suite.T(), &json.UnmarshalTypeError{}, err)
}

func (suite *APIServiceSuite) TestGetMessages_Success() {
	fixture := []map[string]any{
		{
			"id":        "mess001",
			"accountId": "acc001",
			"msgid":     "smth001",
			"from": map[string]string{
				"address": "foo@bar.baz",
				"name":    "Foo Bar",
			},
			"to": []map[string]string{
				{
					"address": "bar@baz.foo",
					"name":    "Bar Baz",
				},
			},
			"subject":        "Foo",
			"intro":          "Foo Bar Baz",
			"seen":           true,
			"isDeleted":      false,
			"hasAttachments": false,
			"size":           200,
			"downloadUrl":    "url",
			"createdAt":      "2022-10-06T13:13:16.857Z",
			"updatedAt":      "2022-10-06T13:13:16.857Z",
		},
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusOK, fixture)
	httpmock.RegisterResponder(http.MethodGet, "/messages", responder)

	response, err := suite.Service.GetMessages()

	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), dto.Messages{}, response)
	assert.Equal(suite.T(), 1, len(response))
	assert.Equal(suite.T(), "mess001", response[0].ID)
	assert.False(suite.T(), response[0].IsDeleted)
}

func (suite *APIServiceSuite) TestGetMessages_InternalServerError() {
	fixture := map[string]any{
		"code":    500,
		"message": "Internal Server Error",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, fixture)
	httpmock.RegisterResponder(http.MethodGet, "/messages", responder)

	response, err := suite.Service.GetMessages()

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.IsType(suite.T(), &errors.HTTPError{}, err)
	assert.Equal(suite.T(), "Internal Server Error [500]", err.Error())
}

func (suite *APIServiceSuite) TestGetMessages_InvalidResponse() {
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, "invalid data")
	httpmock.RegisterResponder(http.MethodGet, "/messages", responder)

	_, err := suite.Service.GetMessages()

	assert.IsType(suite.T(), &json.UnmarshalTypeError{}, err)
}

func (suite *APIServiceSuite) TestGetMessage_Success() {
	fixture := map[string]any{
		"id":        "mess001",
		"accountId": "acc001",
		"msgid":     "smth001",
		"from": map[string]string{
			"address": "foo@bar.baz",
			"name":    "Foo Bar",
		},
		"to": []map[string]string{
			{
				"address": "bar@baz.foo",
				"name":    "Bar Baz",
			},
		},
		"cc": []map[string]string{
			{
				"address": "bar@baz.foo",
				"name":    "Bar Baz",
			},
		},
		"bcc": []map[string]string{
			{
				"address": "bar@baz.foo",
				"name":    "Bar Baz",
			},
		},
		"subject":        "string",
		"seen":           true,
		"flagged":        true,
		"isDeleted":      true,
		"verifications":  []string{},
		"retention":      true,
		"retentionDate":  "2022-10-06T13:20:07.567Z",
		"text":           "foo",
		"html":           []string{},
		"hasAttachments": true,
		"attachments": []map[string]any{
			{
				"id":               "ATTACH000001",
				"filename":         "happy.png",
				"contentType":      "image/png",
				"disposition":      "attachment",
				"transferEncoding": "base64",
				"related":          false,
				"size":             128,
				"downloadUrl":      "/messages/id/attachment/ATTACH000001",
			},
		},
		"size":        128,
		"downloadUrl": "foo",
		"createdAt":   "2022-10-06T13:20:07.567Z",
		"updatedAt":   "2022-10-06T13:20:07.567Z",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusOK, fixture)
	httpmock.RegisterResponder(http.MethodGet, `=~/messages/\w+\z`, responder)

	response, err := suite.Service.GetMessage("mess001")

	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), &dto.Message{}, response)
	assert.Equal(suite.T(), "mess001", response.ID)
	assert.Equal(suite.T(), "ATTACH000001", response.Attachments[0].ID)
}

func (suite *APIServiceSuite) TestGetMessage_NotFound() {
	fixture := map[string]any{
		"type":   "https://tools.ietf.org/html/rfc2616#section-10",
		"title":  "An error occurred",
		"detail": "Not Found",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusNotFound, fixture)
	httpmock.RegisterResponder(http.MethodGet, `=~/messages/\w+\z`, responder)

	_, err := suite.Service.GetMessage("mess001")

	assert.NotNil(suite.T(), err)
	assert.IsType(suite.T(), &errors.HTTPError{}, err)
	assert.Equal(suite.T(), "Not Found [404]", err.Error())
}

func (suite *APIServiceSuite) TestGetMessage_InvalidResponse() {
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, "invalid data")
	httpmock.RegisterResponder(http.MethodGet, `=~/messages/\w+\z`, responder)

	_, err := suite.Service.GetMessage("mess001")

	assert.IsType(suite.T(), &json.UnmarshalTypeError{}, err)
}

func (suite *APIServiceSuite) TestRemoveMessage_Success() {
	fixture := map[string]any{}
	responder, _ := httpmock.NewJsonResponder(http.StatusNoContent, fixture)
	httpmock.RegisterResponder(http.MethodDelete, `=~/messages/\w+\z`, responder)

	err := suite.Service.RemoveMessage("mess001")

	assert.Nil(suite.T(), err)
}

func (suite *APIServiceSuite) TestRemoveMessage_InternalServerError() {
	fixture := map[string]any{
		"type":   "https://tools.ietf.org/html/rfc2616#section-10",
		"title":  "An error occurred",
		"detail": "Internal Server Error",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, fixture)
	httpmock.RegisterResponder(http.MethodDelete, `=~/messages/\w+\z`, responder)

	err := suite.Service.RemoveMessage("mess001")

	assert.NotNil(suite.T(), err)
	assert.IsType(suite.T(), &errors.HTTPError{}, err)
	assert.Equal(suite.T(), "Internal Server Error [500]", err.Error())
}

func (suite *APIServiceSuite) TestRemoveMessage_InvalidResponse() {
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, "invalid data")
	httpmock.RegisterResponder(http.MethodDelete, `=~/messages/\w+\z`, responder)

	err := suite.Service.RemoveMessage("mess001")

	assert.IsType(suite.T(), &json.UnmarshalTypeError{}, err)
}

func (suite *APIServiceSuite) TestUpdateMessage_Success() {
	fixture := map[string]any{}
	responder, _ := httpmock.NewJsonResponder(http.StatusOK, fixture)
	httpmock.RegisterResponder(http.MethodPatch, `=~/messages/\w+\z`, responder)

	err := suite.Service.UpdateMessage("mess001", dto.MessageWrite{Seen: true})

	assert.Nil(suite.T(), err)
}

func (suite *APIServiceSuite) TestUpdateMessage_InvalidInput() {
	fixture := map[string]any{
		"type":   "https://tools.ietf.org/html/rfc2616#section-10",
		"title":  "An error occurred",
		"detail": "seen: This value is required.",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusBadRequest, fixture)
	httpmock.RegisterResponder(http.MethodPatch, `=~/messages/\w+\z`, responder)

	err := suite.Service.UpdateMessage("mess001", dto.MessageWrite{})

	assert.NotNil(suite.T(), err)
	assert.IsType(suite.T(), &errors.HTTPError{}, err)
	assert.Equal(suite.T(), "seen: This value is required. [400]", err.Error())
}

func (suite *APIServiceSuite) TestUpdateMessage_InvalidResponse() {
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, "invalid data")
	httpmock.RegisterResponder(http.MethodPatch, `=~/messages/\w+\z`, responder)

	err := suite.Service.UpdateMessage("mess001", dto.MessageWrite{})

	assert.IsType(suite.T(), &json.UnmarshalTypeError{}, err)
}

func (suite *APIServiceSuite) TestDownloadMessageAttachment_Success() {
	responder := httpmock.NewBytesResponder(http.StatusOK, []byte{})
	httpmock.RegisterResponder(http.MethodGet, `=~/messages/\w+/attachment/\w+\z`, responder)

	path := "/tmp/attachments/ATTACH000001"
	err := suite.Service.DownloadMessageAttachment("mess1", "ATTACH000001", path)

	assert.Nil(suite.T(), err)

	info, err := os.Lstat(path)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "ATTACH000001", info.Name())
}

func (suite *APIServiceSuite) TestDownloadMessageAttachment_NotFound() {
	responder := httpmock.NewBytesResponder(http.StatusNotFound, []byte{})
	httpmock.RegisterResponder(http.MethodGet, `=~/messages/\w+/attachment/\w+\z`, responder)

	err := suite.Service.DownloadMessageAttachment("mess1", "ATTACH000001", "/tmp/bar")

	assert.NotNil(suite.T(), err)
	assert.IsType(suite.T(), &errors.HTTPError{}, err)
	assert.Equal(suite.T(), " [404]", err.Error())
}

func (suite *APIServiceSuite) TestDownloadMessageAttachment_InvalidPath() {
	responder := httpmock.NewBytesResponder(http.StatusNotFound, []byte{})
	httpmock.RegisterResponder(http.MethodGet, `=~/messages/\w+/attachment/\w+\z`, responder)

	err := suite.Service.DownloadMessageAttachment("mess1", "ATTACH000001", "/invalid/path")

	assert.IsType(suite.T(), &fs.PathError{}, err)
}

func (suite *APIServiceSuite) TestGetSource_Success() {
	fixture := map[string]any{
		"id":          "src001",
		"downloadUrl": "foo",
		"data":        "bar",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusOK, fixture)
	httpmock.RegisterResponder(http.MethodGet, `=~/sources/\w+\z`, responder)

	response, err := suite.Service.GetSource("src001")

	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), &dto.Source{}, response)
	assert.Equal(suite.T(), "src001", response.ID)
	assert.Equal(suite.T(), "bar", response.Data)
}

func (suite *APIServiceSuite) TestGetSource_NotFound() {
	fixture := map[string]any{
		"type":   "https://tools.ietf.org/html/rfc2616#section-10",
		"title":  "An error occurred",
		"detail": "Not Found",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusNotFound, fixture)
	httpmock.RegisterResponder(http.MethodGet, `=~/sources/\w+\z`, responder)

	_, err := suite.Service.GetSource("src001")

	assert.NotNil(suite.T(), err)
	assert.IsType(suite.T(), &errors.HTTPError{}, err)
	assert.Equal(suite.T(), "Not Found [404]", err.Error())
}

func (suite *APIServiceSuite) TestGetSource_InvalidResponse() {
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, "invalid data")
	httpmock.RegisterResponder(http.MethodGet, `=~/sources/\w+\z`, responder)

	_, err := suite.Service.GetSource("src001")

	assert.IsType(suite.T(), &json.UnmarshalTypeError{}, err)
}

func (suite *APIServiceSuite) TestGetToken_Success() {
	fixture := map[string]any{
		"id":    "acc001",
		"token": "foo",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusOK, fixture)
	httpmock.RegisterResponder(http.MethodPost, "/token", responder)

	body := dto.Credentials{
		Address:  "foo@bar.baz",
		Password: "Pa$$w0rd",
	}
	response, err := suite.Service.GetToken(body)

	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), &dto.Token{}, response)
	assert.Equal(suite.T(), "acc001", response.ID)
	assert.Equal(suite.T(), "foo", response.Token)
}

func (suite *APIServiceSuite) TestGetToken_InvalidCredentials() {
	fixture := map[string]any{
		"code":    401,
		"message": "Invalid credentials.",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusUnauthorized, fixture)
	httpmock.RegisterResponder(http.MethodPost, "/token", responder)

	body := dto.Credentials{
		Address:  "invalid@foo.bar",
		Password: "Pa$$w0rd",
	}
	response, err := suite.Service.GetToken(body)

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.IsType(suite.T(), &errors.HTTPError{}, err)
	assert.Equal(suite.T(), "Invalid credentials. [401]", err.Error())
}

func (suite *APIServiceSuite) TestGetToken_InvalidResponse() {
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, "invalid data")
	httpmock.RegisterResponder(http.MethodPost, "/token", responder)

	_, err := suite.Service.GetToken(dto.Credentials{})

	assert.IsType(suite.T(), &json.UnmarshalTypeError{}, err)
}

func (suite *APIServiceSuite) TestSetToken() {
	token := &dto.Token{
		ID:    "foo",
		Token: "bar",
	}
	suite.Service.SetToken(token)

	assert.Equal(suite.T(), "bar", suite.Service.client.Token)
}
