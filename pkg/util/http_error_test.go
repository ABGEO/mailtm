package util

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HTTPErrorSuite struct {
	suite.Suite

	Client *resty.Client
}

func (suite *HTTPErrorSuite) SetupSuite() {
	suite.Client = resty.New()

	httpmock.ActivateNonDefault(suite.Client.GetClient())
}

func (suite *HTTPErrorSuite) TearDownTest() {
	httpmock.Reset()
}

func (suite *HTTPErrorSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func TestHTTPErrorSuite(t *testing.T) {
	suite.Run(t, new(HTTPErrorSuite))
}

func (suite *HTTPErrorSuite) TestNewHTTPError_ResponseNotAHTTPError() {
	fixture := map[string]any{}
	responder, _ := httpmock.NewJsonResponder(http.StatusNotFound, fixture)
	httpmock.RegisterResponder(http.MethodGet, "/foo", responder)

	resp, _ := suite.Client.R().Get("/foo")
	httpError := NewHTTPError(resp)

	assert.Equal(suite.T(), " [404]", httpError.Error())
}

func (suite *HTTPErrorSuite) TestNewHTTPError_ResponseIsAPrimitiveError() {
	fixture := map[string]any{
		"code":    "404",
		"message": "Not Found",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusNotFound, fixture)
	httpmock.RegisterResponder(http.MethodGet, "/foo", responder)

	resp, _ := suite.Client.R().SetError(&HTTPError{}).Get("/foo")
	httpError := NewHTTPError(resp)

	assert.Equal(suite.T(), "Not Found [404]", httpError.Error())
}

func (suite *HTTPErrorSuite) TestNewHTTPError_NormalError() {
	fixture := map[string]any{
		"type":   "https://tools.ietf.org/html/rfc2616#section-10",
		"title":  "An error occurred",
		"detail": "address: This value is not a valid email address.",
	}
	responder, _ := httpmock.NewJsonResponder(http.StatusBadRequest, fixture)
	httpmock.RegisterResponder(http.MethodGet, "/foo", responder)

	resp, _ := suite.Client.R().SetError(&HTTPError{}).Get("/foo")
	httpError := NewHTTPError(resp)

	assert.Equal(suite.T(), "address: This value is not a valid email address. [400]", httpError.Error())
}
