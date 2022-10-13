package auth

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/errors"
	"github.com/abgeo/mailtm/test"
	"github.com/abgeo/mailtm/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LoginCmdSuite struct {
	test.BaseCMDSuite

	APIServiceMock *mocks.APIServiceInterface
}

func TestLoginCmdSuite(t *testing.T) {
	suite.Run(t, new(LoginCmdSuite))
}

func (suite *LoginCmdSuite) SetupSuite() {
	suite.Buffer = bytes.NewBufferString("")
	suite.APIServiceMock = mocks.NewAPIServiceInterface(suite.T())
	suite.CmdOptions = command.Options{
		Writer:     suite.Buffer,
		APIService: suite.APIServiceMock,
	}

	errors.SetErrorHandler(func(msg interface{}, code int) {
		fmt.Fprintln(suite.Buffer, "Error:", msg)
	})
}

func (suite *LoginCmdSuite) TearDownTest() {
	suite.Buffer.Reset()
}

func (suite *LoginCmdSuite) TestNewCmdLogin_Success() {
	getTokenFixture := &dto.Token{
		ID:    "acc001",
		Token: "foo",
	}
	suite.APIServiceMock.On("GetToken", mock.Anything).Return(getTokenFixture, nil).Times(1)

	cmd := NewCmdLogin(suite.CmdOptions)
	cmd.SetArgs([]string{"foo@bar.baz", "pass"})
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Contains(suite.T(), cmdOutput, "User foo@bar.baz has been authenticated successfully.")
}

func (suite *LoginCmdSuite) TestNewCmdLogin_InvalidCredentials() {
	fixtureError := &errors.HTTPError{
		Code:   401,
		Detail: "Invalid Credentials.",
	}
	suite.APIServiceMock.On("GetToken", mock.Anything).Return(nil, fixtureError).Times(1)

	cmd := NewCmdLogin(suite.CmdOptions)
	cmd.SetArgs([]string{"foo@bar.baz", "pass"})
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Contains(suite.T(), cmdOutput, "Invalid Credentials. [401]")
}
