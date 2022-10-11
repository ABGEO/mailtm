package account

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

type AccountCMDSuite struct {
	test.BaseCMDSuite

	APIServiceMock *mocks.APIServiceInterface
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

	suite.APIServiceMock = mocks.NewAPIServiceInterface(suite.T())
	suite.CmdOptions = command.Options{
		Version:    appVersion,
		APIService: suite.APIServiceMock,
	}
}

func (suite *AccountCMDSuite) TestAccountRootCMD() {
	suite.APIServiceMock.On("SetToken", mock.Anything)

	cmd := NewCmd(suite.CmdOptions)

	assert.Contains(suite.T(), suite.GetCommandOutput(cmd), "Usage:\n  account [command]")
}
