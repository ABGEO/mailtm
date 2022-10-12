package auth

import (
	"testing"

	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthCmdSuite struct {
	test.BaseCMDSuite
}

func TestAuthCmdSuite(t *testing.T) {
	suite.Run(t, new(AuthCmdSuite))
}

func (suite *AuthCmdSuite) SetupSuite() {
	suite.CmdOptions = command.Options{}
}

func (suite *AuthCmdSuite) TestAccountRootCmd() {
	cmd := NewCmd(suite.CmdOptions)

	assert.Contains(suite.T(), suite.GetCommandOutput(cmd), "Usage:\n  auth [command]")
}
