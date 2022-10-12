package message

import (
	"testing"

	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/test"
	"github.com/abgeo/mailtm/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MessageCmdSuite struct {
	test.BaseCMDSuite

	APIServiceMock *mocks.APIServiceInterface
}

func (suite *MessageCmdSuite) SetupSuite() {
	suite.APIServiceMock = mocks.NewAPIServiceInterface(suite.T())
	suite.CmdOptions = command.Options{
		APIService: suite.APIServiceMock,
	}

	suite.APIServiceMock.On("SetToken", mock.Anything)
}

func TestMessageCmdSuite(t *testing.T) {
	suite.Run(t, new(MessageCmdSuite))
}

func (suite *MessageCmdSuite) TestAccountRootCmd() {
	cmd := NewCmd(suite.CmdOptions)

	assert.Contains(suite.T(), suite.GetCommandOutput(cmd), "Usage:\n  message [command]")
}
