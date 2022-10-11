package account

import (
	"testing"

	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/test"
	"github.com/abgeo/mailtm/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AccountCmdSuite struct {
	test.BaseCMDSuite

	APIServiceMock *mocks.APIServiceInterface
}

func TestAccountCmdSuite(t *testing.T) {
	suite.Run(t, new(AccountCmdSuite))
}

func (suite *AccountCmdSuite) SetupSuite() {
	suite.APIServiceMock = mocks.NewAPIServiceInterface(suite.T())
	suite.CmdOptions = command.Options{
		APIService: suite.APIServiceMock,
	}
}

func (suite *AccountCmdSuite) TestAccountRootCmd() {
	suite.APIServiceMock.On("SetToken", mock.Anything)

	cmd := NewCmd(suite.CmdOptions)

	assert.Contains(suite.T(), suite.GetCommandOutput(cmd), "Usage:\n  account [command]")
}
