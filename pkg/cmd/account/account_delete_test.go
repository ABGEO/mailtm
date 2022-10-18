package account

import (
	"bytes"
	"fmt"
	"testing"

	"atomicgo.dev/keyboard"
	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/errors"
	"github.com/abgeo/mailtm/test"
	"github.com/abgeo/mailtm/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DeleteCmdSuite struct {
	test.BaseCMDSuite

	APIServiceMock *mocks.APIServiceInterface
}

func (suite *DeleteCmdSuite) SetupSuite() {
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

func (suite *DeleteCmdSuite) TearDownTest() {
	suite.Buffer.Reset()
}

func TestDeleteCmdSuite(t *testing.T) {
	suite.Run(t, new(DeleteCmdSuite))
}

func (suite *DeleteCmdSuite) TestNewCmdDelete_Success() {
	suite.APIServiceMock.On("RemoveAccount", mock.Anything).Return(nil).Times(1)

	go func() {
		err := keyboard.SimulateKeyPress("y")
		if err != nil {
			assert.Error(suite.T(), err)
		}
	}()

	cmd := NewCmdDelete(suite.CmdOptions)
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Equal(suite.T(), "Account has been deleted successfully.\n", cmdOutput)
}

func (suite *DeleteCmdSuite) TestNewCmdDelete_DoNotDelete() {
	go func() {
		err := keyboard.SimulateKeyPress("n")
		if err != nil {
			assert.Error(suite.T(), err)
		}
	}()

	cmd := NewCmdDelete(suite.CmdOptions)
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Equal(suite.T(), "", cmdOutput)
}

func (suite *DeleteCmdSuite) TestNewCmdDelete_AccessDenied() {
	fixtureError := &errors.HTTPError{
		Code:   401,
		Detail: "Access Denied",
	}
	suite.APIServiceMock.On("RemoveAccount", mock.Anything).Return(fixtureError).Times(1)

	go func() {
		err := keyboard.SimulateKeyPress("y")
		if err != nil {
			assert.Error(suite.T(), err)
		}
	}()

	cmd := NewCmdDelete(suite.CmdOptions)
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Equal(suite.T(), cmdOutput, "Error: Access Denied [401]\n")
}
