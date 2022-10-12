package auth

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/errors"
	"github.com/abgeo/mailtm/test"
	"github.com/abgeo/mailtm/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RandomCmdSuite struct {
	test.BaseCMDSuite

	APIServiceMock *mocks.APIServiceInterface
}

func TestRandomCmdSuite(t *testing.T) {
	suite.Run(t, new(RandomCmdSuite))
}

func (suite *RandomCmdSuite) SetupSuite() {
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

func (suite *RandomCmdSuite) TearDownTest() {
	suite.Buffer.Reset()
}

func (suite *RandomCmdSuite) TestNewCmdRandom_Success() {
	getDomainsFixture := []dto.Domain{
		{
			ID:        "dom001",
			Domain:    "bar.com",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	createAccountFixture := &dto.Account{
		ID:         "acc001",
		Address:    "foo@bar.baz",
		Quota:      40000000,
		Used:       0,
		IsDisabled: false,
		IsDeleted:  false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	getTokenFixture := &dto.Token{
		ID:    "acc001",
		Token: "foo",
	}

	suite.APIServiceMock.On("GetDomains").Return(getDomainsFixture, nil).Times(1)
	suite.APIServiceMock.On("CreateAccount", mock.Anything).Return(createAccountFixture, nil).Times(1)
	suite.APIServiceMock.On("GetToken", mock.Anything).Return(getTokenFixture, nil).Times(1)

	cmd := NewCmdRandom(suite.CmdOptions)
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Contains(suite.T(), cmdOutput, "New random account has been created")
	assert.Regexp(suite.T(), fmt.Sprintf(`Email[ ]*:[ ]*%s`, createAccountFixture.Address), cmdOutput)
	assert.Regexp(suite.T(), `Password[ ]*:[ ]*\w*`, cmdOutput)
}

func (suite *RandomCmdSuite) TestNewCmdRandom_ErrorFetchingDomains() {
	fixtureError := &errors.HTTPError{
		Code:   500,
		Detail: "Internal Server Error",
	}
	suite.APIServiceMock.On("GetDomains").Return(nil, fixtureError).Times(1)

	cmd := NewCmdRandom(suite.CmdOptions)
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Contains(suite.T(), cmdOutput, "Error: Internal Server Error [500]")
}

func (suite *RandomCmdSuite) TestNewCmdRandom_ErrorCreatingAccount() {
	getDomainsFixture := []dto.Domain{
		{
			ID:        "dom001",
			Domain:    "bar.com",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	createAccountFixture := &dto.Account{}
	createAccountFixtureError := &errors.HTTPError{
		Code:   400,
		Detail: "User Already Exists",
	}

	suite.APIServiceMock.On("GetDomains").Return(getDomainsFixture, nil).Times(1)
	suite.APIServiceMock.On("CreateAccount", mock.Anything).
		Return(createAccountFixture, createAccountFixtureError).
		Times(1)

	cmd := NewCmdRandom(suite.CmdOptions)
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Contains(suite.T(), cmdOutput, "Error: User Already Exists [400]")
}

func (suite *RandomCmdSuite) TestNewCmdRandom_ErrorFetchingToken() {
	getDomainsFixture := []dto.Domain{
		{
			ID:        "dom001",
			Domain:    "bar.com",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	createAccountFixture := &dto.Account{
		ID:         "acc001",
		Address:    "foo@bar.baz",
		Quota:      40000000,
		Used:       0,
		IsDisabled: false,
		IsDeleted:  false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	getTokenFixture := &dto.Token{}
	getTokenFixtureError := &errors.HTTPError{
		Code:   401,
		Detail: "Access Denied",
	}

	suite.APIServiceMock.On("GetDomains").Return(getDomainsFixture, nil).Times(1)
	suite.APIServiceMock.On("CreateAccount", mock.Anything).Return(createAccountFixture, nil).Times(1)
	suite.APIServiceMock.On("GetToken", mock.Anything).Return(getTokenFixture, getTokenFixtureError).Times(1)

	cmd := NewCmdRandom(suite.CmdOptions)
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Contains(suite.T(), cmdOutput, "Error: Access Denied [401]")
}
