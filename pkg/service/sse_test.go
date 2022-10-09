package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/util"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var errTimeout = errors.New("timeout")

type SSEServiceSuite struct {
	suite.Suite

	Service *SSEService
}

func (suite *SSEServiceSuite) SetupSuite() {
	version := util.Version{
		Number: "test",
	}
	authConf := configs.AuthConfig{
		ID:    "acc001",
		Email: "foo@bar.baz",
		Token: "foo",
	}
	suite.Service = NewSSEService(version, authConf)

	httpmock.ActivateNonDefault(suite.Service.client.Connection)
}

func (suite *SSEServiceSuite) TearDownTest() {
	httpmock.Reset()
}

func (suite *SSEServiceSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func TestSSEServiceSuite(t *testing.T) {
	suite.Run(t, new(SSEServiceSuite))
}

func (suite *SSEServiceSuite) TestUserAgent() {
	assert.Equal(suite.T(), "abgeo-mailtm/test", suite.Service.client.Headers["User-Agent"])
}

func (suite *SSEServiceSuite) TestSubscribeMessages_Success() {
	//nolint:errchkjson
	fixture, _ := json.Marshal(map[string]any{
		"@context":  "/contexts/Message",
		"@id":       "/messages/mess001",
		"@type":     "Message",
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
		"seen":           false,
		"isDeleted":      false,
		"hasAttachments": false,
		"size":           200,
		"downloadUrl":    "url",
		"createdAt":      "2022-10-06T13:13:16.857Z",
		"updatedAt":      "2022-10-06T13:13:16.857Z",
	})

	responder := httpmock.NewStringResponder(
		http.StatusOK,
		fmt.Sprintf("id: urn:uuid:ba1055ae-452b-4e21-a77a-5c6fcc200178\ndata: %s\n\n", fixture),
	)
	httpmock.RegisterResponder(http.MethodGet, `=~/\.well-known/mercure\?topic=/accounts/\w+\z`, responder)

	messages := make(chan dto.MessagesItem)

	go func() {
		_ = suite.Service.SubscribeMessages("acc001", func(message dto.MessagesItem) {
			messages <- message
		})
	}()

	message, err := suite.waitForMessage(messages, 1*time.Second)

	assert.NotNil(suite.T(), message)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "mess001", message.ID)
	assert.Equal(suite.T(), false, message.Seen)
}

func (suite *SSEServiceSuite) waitForMessage(ch chan dto.MessagesItem, duration time.Duration) (
	message dto.MessagesItem,
	err error,
) {
	select {
	case message = <-ch:
		return message, nil
	case <-time.After(duration):
		return message, errTimeout
	}
}
