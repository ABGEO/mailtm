package message

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/errors"
	"github.com/abgeo/mailtm/test"
	"github.com/abgeo/mailtm/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DownloadAttachmentCmdSuite struct {
	test.BaseCMDSuite

	APIServiceMock *mocks.APIServiceInterface
}

func (suite *DownloadAttachmentCmdSuite) SetupSuite() {
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

func (suite *DownloadAttachmentCmdSuite) TearDownTest() {
	suite.Buffer.Reset()
}

func TestDownloadAttachmentCmdSuite(t *testing.T) {
	suite.Run(t, new(DownloadAttachmentCmdSuite))
}

func (suite *DownloadAttachmentCmdSuite) TestNewCmdDownloadAttachment_Success() {
	getMessageFixture := &dto.Message{
		Attachments: []dto.MessageAttachment{{
			ID:       "ATTACH000001",
			Filename: "happy.png",
		}},
	}

	suite.APIServiceMock.On("GetMessage", mock.Anything).Return(getMessageFixture, nil).Times(1)
	suite.APIServiceMock.On("DownloadMessageAttachment", "mess001", "ATTACH000001", "./attachments/happy.png").
		Return(nil).
		Times(1)

	cmd := NewCmdDownloadAttachment(suite.CmdOptions)
	cmd.SetArgs([]string{"mess001", "ATTACH000001", "--dir", "./attachments"})
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Equal(suite.T(), "", cmdOutput)
}

func (suite *DownloadAttachmentCmdSuite) TestNewCmdDownloadAttachment_MessageNotFound() {
	getMessageFixture := &dto.Message{}
	getMessageFixtureError := &errors.HTTPError{
		Code:   404,
		Detail: "Not Found",
	}

	suite.APIServiceMock.On("GetMessage", mock.Anything).Return(getMessageFixture, getMessageFixtureError).Times(1)

	cmd := NewCmdDownloadAttachment(suite.CmdOptions)
	cmd.SetArgs([]string{"mess001", "ATTACH000001"})
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Contains(suite.T(), cmdOutput, "Error: Not Found [404]")
}

func (suite *DownloadAttachmentCmdSuite) TestNewCmdDownloadAttachment_AttachmentNotFound() {
	getMessageFixture := &dto.Message{}

	suite.APIServiceMock.On("GetMessage", mock.Anything).Return(getMessageFixture, nil).Times(1)

	cmd := NewCmdDownloadAttachment(suite.CmdOptions)
	cmd.SetArgs([]string{"mess001", "ATTACH000001"})
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Contains(suite.T(), cmdOutput, "Error: attachment not found")
}
