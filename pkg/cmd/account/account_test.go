package account

import (
	"testing"

	"github.com/abgeo/mailtm/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AccountCMDSuite struct {
	test.BaseCMDSuite
}

func TestAccountCMDSuite(t *testing.T) {
	suite.Run(t, new(AccountCMDSuite))
}

func (suite *AccountCMDSuite) TestAccountRootCMD() {
	cmd := NewCmd(suite.CmdOptions)

	assert.Contains(suite.T(), suite.GetCommandOutput(cmd), "Usage:\n  account [command]")
}
