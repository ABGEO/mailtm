package account

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

type CurrentCmdSuite struct {
	test.BaseCMDSuite

	APIServiceMock *mocks.APIServiceInterface
}

func TestCurrentCmdSuite(t *testing.T) {
	suite.Run(t, new(CurrentCmdSuite))
}

func (suite *CurrentCmdSuite) SetupSuite() {
	suite.Buffer = bytes.NewBufferString("")
	suite.APIServiceMock = mocks.NewAPIServiceInterface(suite.T())
	suite.CmdOptions = command.Options{
		Writer:     suite.Buffer,
		APIService: suite.APIServiceMock,
	}

	errors.SetErrorHandler(func(msg interface{}, code int) {
		fmt.Fprintln(suite.Buffer, "Error:", msg)
	})

	suite.APIServiceMock.On("SetToken", mock.Anything)
}

func (suite *CurrentCmdSuite) TearDownTest() {
	suite.Buffer.Reset()
}

func (suite *CurrentCmdSuite) TestNewCmdCurrent_Success() {
	fixture := &dto.Account{
		ID:        "acc001",
		Address:   "foo@bar.baz",
		Quota:     40000000,
		Used:      10000000,
		CreatedAt: time.Now(),
	}
	suite.APIServiceMock.On("GetCurrentAccount").Return(fixture, nil).Times(1)

	cmd := NewCmdCurrent(suite.CmdOptions)
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Regexp(suite.T(), fmt.Sprintf(`ID[ ]*:[ ]*%s`, fixture.ID), cmdOutput)
	assert.Regexp(suite.T(), fmt.Sprintf(`Email[ ]*:[ ]*%s`, fixture.Address), cmdOutput)
	assert.Regexp(suite.T(), `Usage[ ]*:[ ]*9.54 MB / 38.15 MB`, cmdOutput)
	assert.Regexp(
		suite.T(),
		fmt.Sprintf(`Created At[ ]*:[ ]*%s`, fixture.CreatedAt.Format("02 January 2006 15:04:05")),
		cmdOutput,
	)
}

func (suite *CurrentCmdSuite) TestNewCmdCurrent_AccessDenied() {
	fixture := &dto.Account{}
	fixtureError := &errors.HTTPError{
		Code:   401,
		Detail: "Access Denied",
	}
	suite.APIServiceMock.On("GetCurrentAccount").Return(fixture, fixtureError).Times(1)

	cmd := NewCmdCurrent(suite.CmdOptions)
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Contains(suite.T(), cmdOutput, "Error: Access Denied [401]")
}
