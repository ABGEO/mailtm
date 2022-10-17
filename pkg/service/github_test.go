package service

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GitHubServiceSuite struct {
	suite.Suite

	Service *GitHubService
}

func (suite *GitHubServiceSuite) SetupSuite() {
	suite.Service = NewGitHubService()

	httpmock.ActivateNonDefault(suite.Service.client.GetClient())
}

func (suite *GitHubServiceSuite) TearDownTest() {
	httpmock.Reset()
}

func (suite *GitHubServiceSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func TestGitHubServiceSuite(t *testing.T) {
	suite.Run(t, new(GitHubServiceSuite))
}

func (suite *GitHubServiceSuite) TestGetLatestRelease_Success() {
	fixture := map[string]interface{}{
		"id":           100,
		"html_url":     "https://foo.bar/v1.0.0",
		"name":         "v1.0.0",
		"tag_name":     "v1.0.0",
		"created_at":   time.Now(),
		"published_at": time.Now(),
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusOK, fixture)
	httpmock.RegisterResponder(http.MethodGet, `=~/repos/\w+/\w+/releases/latest\z`, responder)

	response, err := suite.Service.GetLatestRelease("foo", "bar")

	assert.Nil(suite.T(), err)
	assert.IsType(suite.T(), &dto.Release{}, response)
	assert.Equal(suite.T(), fixture["id"], response.ID)
	assert.Equal(suite.T(), fixture["name"], response.Name)
}

func (suite *GitHubServiceSuite) TestGetLatestRelease_InvalidResponse() {
	responder, _ := httpmock.NewJsonResponder(http.StatusInternalServerError, "invalid data")
	httpmock.RegisterResponder(http.MethodGet, `=~/repos/\w+/\w+/releases/latest\z`, responder)

	_, err := suite.Service.GetLatestRelease("foo", "bar")

	assert.IsType(suite.T(), &json.UnmarshalTypeError{}, err)
}
