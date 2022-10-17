package cmd

import (
	"testing"

	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/types"
	"github.com/abgeo/mailtm/test"
	"github.com/abgeo/mailtm/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CMDSuite struct {
	test.BaseCMDSuite

	APIServiceMock *mocks.APIServiceInterface
}

func (suite *CMDSuite) SetupSuite() {
	suite.APIServiceMock = mocks.NewAPIServiceInterface(suite.T())
	appVersion := types.Version{
		Number: "foo",
		Date:   "bar",
		Commit: "baz",
	}

	suite.CmdOptions = command.Options{
		Version:    appVersion,
		APIService: suite.APIServiceMock,
	}

	suite.APIServiceMock.On("SetToken", mock.Anything)
}

func TestCMDSuite(t *testing.T) {
	suite.Run(t, new(CMDSuite))
}

func (suite *CMDSuite) TestNewCmd_WithoutFlags() {
	cmd := NewCmd(suite.CmdOptions)

	assert.Contains(suite.T(), suite.GetCommandOutput(cmd), "CLI client for Mail.tm disposable mail service")
}

func (suite *CMDSuite) TestNewCmd_WithVersionFlag() {
	cmd := NewCmd(suite.CmdOptions)
	cmd.SetArgs([]string{"-v"})

	assert.Contains(suite.T(), suite.GetCommandOutput(cmd), "mailtm version foo (bar)\nbaz")
}
