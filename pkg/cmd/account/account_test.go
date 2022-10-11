package account

import (
	"testing"

	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/abgeo/mailtm/pkg/types"
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

func (suite *AccountCMDSuite) SetupSuite() {
	appVersion := types.Version{
		Number: "foo",
		Date:   "bar",
		Commit: "baz",
	}

	suite.CmdOptions = command.Options{
		Version: appVersion,
		// @TODO: Replace with mock.
		APIService: service.NewAPIService(appVersion),
	}
}

func (suite *AccountCMDSuite) TestAccountRootCMD() {
	cmd := NewCmd(suite.CmdOptions)

	assert.Contains(suite.T(), suite.GetCommandOutput(cmd), "Usage:\n  account [command]")
}
