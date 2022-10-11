package cmd

import (
	"testing"

	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/abgeo/mailtm/pkg/types"
	"github.com/abgeo/mailtm/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CMDSuite struct {
	test.BaseCMDSuite
}

func (suite *CMDSuite) SetupSuite() {
	appVersion := types.Version{
		Number: "foo",
		Date:   "bar",
		Commit: "baz",
	}

	suite.CmdOptions = command.Options{
		Version: appVersion,
		// @TODO: Replace with mocks.
		APIService: service.NewAPIService(appVersion),
	}
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
