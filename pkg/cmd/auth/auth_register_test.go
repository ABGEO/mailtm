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

type RegisterCmdSuite struct {
	test.BaseCMDSuite

	APIServiceMock *mocks.APIServiceInterface
}

func TestRegisterCmdSuite(t *testing.T) {
	suite.Run(t, new(RegisterCmdSuite))
}

func (suite *RegisterCmdSuite) SetupSuite() {
	suite.Buffer = bytes.NewBufferString("")
	suite.APIServiceMock = mocks.NewAPIServiceInterface(suite.T())
	suite.CmdOptions = command.Options{
		Writer:     suite.Buffer,
		APIService: suite.APIServiceMock,
	}

	errors.SetErrorHandler(func(msg interface{}, code int) {
		fmt.Fprintln(suite.Buffer, "Error:", msg)
		panic(fmt.Sprintf("exited with %d", code))
	})
}

func (suite *RegisterCmdSuite) TearDownTest() {
	suite.Buffer.Reset()
}

func (suite *RegisterCmdSuite) TestNewCmdRegister_WithoutLogin() {
	getDomainsFixture := []dto.Domain{
		{Domain: "bar.baz"},
	}
	createAccountFixture := &dto.Account{
		Address: "foo@bar.baz",
	}

	suite.APIServiceMock.On("GetDomains").Return(getDomainsFixture, nil).Times(1)
	suite.APIServiceMock.On("CreateAccount", mock.Anything).Return(createAccountFixture, nil).Times(1)

	cmd := NewCmdRegister(suite.CmdOptions)
	cmd.SetArgs([]string{"foo", "pass"})
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Equal(suite.T(), "New account foo@bar.baz has been created successfully.\n", cmdOutput)
}

func (suite *RegisterCmdSuite) TestNewCmdRegister_Login() {
	getDomainsFixture := []dto.Domain{
		{Domain: "bar.baz"},
	}
	createAccountFixture := &dto.Account{
		Address: "foo@bar.baz",
	}
	getTokenFixture := &dto.Token{
		ID:    "acc001",
		Token: "foo",
	}

	suite.APIServiceMock.On("GetDomains").Return(getDomainsFixture, nil).Times(1)
	suite.APIServiceMock.On("CreateAccount", mock.Anything).Return(createAccountFixture, nil).Times(1)
	suite.APIServiceMock.On("GetToken", mock.Anything).Return(getTokenFixture, nil).Times(1)

	cmd := NewCmdRegister(suite.CmdOptions)
	cmd.SetArgs([]string{"foo", "pass", "--login"})
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Equal(
		suite.T(),
		"New account foo@bar.baz has been created successfully.\nAuthenticated with new account foo@bar.baz.\n",
		cmdOutput,
	)
}

func (suite *RegisterCmdSuite) TestNewCmdRegister_WithSpecificDomain() {
	getDomainsFixture := []dto.Domain{
		{Domain: "bar.baz"},
		{Domain: "foo.bar"},
	}
	createAccountFixture := &dto.Account{
		Address: "baz@foo.bar",
	}

	suite.APIServiceMock.On("GetDomains").Return(getDomainsFixture, nil).Times(1)
	suite.APIServiceMock.On("CreateAccount", mock.Anything).Return(createAccountFixture, nil).Times(1)

	cmd := NewCmdRegister(suite.CmdOptions)
	cmd.SetArgs([]string{"baz", "pass", "--domain", "foo.bar"})
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Equal(suite.T(), "New account baz@foo.bar has been created successfully.\n", cmdOutput)
}

func (suite *RegisterCmdSuite) TestNewCmdRegister_WithInvalidDomain() {
	getDomainsFixture := []dto.Domain{
		{Domain: "bar.baz"},
		{Domain: "foo.bar"},
	}

	suite.APIServiceMock.On("GetDomains").Return(getDomainsFixture, nil).Times(1)

	cmd := NewCmdRegister(suite.CmdOptions)
	cmd.SetArgs([]string{"baz", "pass", "--domain", "baz.bar"})

	assert.PanicsWithValue(suite.T(), "exited with 1", func() { suite.GetCommandOutput(cmd) })
	assert.Equal(suite.T(),
		"Error: domain baz.bar is not valid. Valid domains are: [bar.baz, foo.bar]\n",
		suite.Buffer.String(),
	)
}

func (suite *RegisterCmdSuite) TestNewCmdRegister_UnableToGetDomains() {
	getDomainsFixtureErr := &errors.HTTPError{
		Code:   500,
		Detail: "Internal Server Error",
	}

	suite.APIServiceMock.On("GetDomains", mock.Anything).Return(nil, getDomainsFixtureErr).Times(1)

	cmd := NewCmdRegister(suite.CmdOptions)
	cmd.SetArgs([]string{"baz", "pass"})

	assert.PanicsWithValue(suite.T(), "exited with 1", func() { suite.GetCommandOutput(cmd) })
	assert.Equal(suite.T(),
		"Error: Internal Server Error [500]\n",
		suite.Buffer.String(),
	)
}

func (suite *RegisterCmdSuite) TestNewCmdRegister_EmailIsInUse() {
	getDomainsFixture := []dto.Domain{
		{Domain: "bar.baz"},
	}
	createAccountFixtureError := &errors.HTTPError{
		Code:   422,
		Detail: "Email is already in use",
	}

	suite.APIServiceMock.On("GetDomains").Return(getDomainsFixture, nil).Times(1)
	suite.APIServiceMock.On("CreateAccount", mock.Anything).Return(nil, createAccountFixtureError).Times(1)

	cmd := NewCmdRegister(suite.CmdOptions)
	cmd.SetArgs([]string{"baz", "pass"})

	assert.PanicsWithValue(suite.T(), "exited with 1", func() { suite.GetCommandOutput(cmd) })
	assert.Equal(suite.T(),
		"Error: Email is already in use [422]\n",
		suite.Buffer.String(),
	)
}

func (suite *RegisterCmdSuite) TestNewCmdRegister_UnableToLogIn() {
	getDomainsFixture := []dto.Domain{
		{Domain: "bar.baz"},
	}
	createAccountFixture := &dto.Account{
		Address: "foo@bar.baz",
	}
	getTokenFixtureError := &errors.HTTPError{
		Code:   401,
		Detail: "Access Denied",
	}

	suite.APIServiceMock.On("GetDomains").Return(getDomainsFixture, nil).Times(1)
	suite.APIServiceMock.On("CreateAccount", mock.Anything).Return(createAccountFixture, nil).Times(1)
	suite.APIServiceMock.On("GetToken", mock.Anything).Return(nil, getTokenFixtureError).Times(1)

	cmd := NewCmdRegister(suite.CmdOptions)
	cmd.SetArgs([]string{"baz", "pass", "-l"})

	assert.PanicsWithValue(suite.T(), "exited with 1", func() { suite.GetCommandOutput(cmd) })
	assert.Contains(suite.T(), suite.Buffer.String(), "Error: Access Denied [401]")
}
