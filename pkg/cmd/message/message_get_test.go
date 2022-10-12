package message

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/errors"
	"github.com/abgeo/mailtm/test"
	"github.com/abgeo/mailtm/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetCmdSuite struct {
	test.BaseCMDSuite

	APIServiceMock *mocks.APIServiceInterface
}

func (suite *GetCmdSuite) SetupSuite() {
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

func (suite *GetCmdSuite) TearDownTest() {
	suite.Buffer.Reset()
}

func TestGetCmdSuite(t *testing.T) {
	suite.Run(t, new(GetCmdSuite))
}

func (suite *GetCmdSuite) TestNewCmdGet_Success() {
	timeNow := time.Now()
	getMessageFixture := &dto.Message{
		Cc: []dto.EmailAddress{{
			Address: "cc@bar.com",
			Name:    "Cc Bar",
		}},
		Bcc: []dto.EmailAddress{{
			Address: "bcc@bar.com",
			Name:    "Bcc Bar",
		}},
		RetentionDate: timeNow.AddDate(0, 0, 7),
		Text:          "Lorem Ipsum is simply dummy text of the printing and typesetting industry.",
	}
	getMessageFixture.ID = "mess001"
	getMessageFixture.Subject = "foo bar baz"
	getMessageFixture.CreatedAt = timeNow
	getMessageFixture.From = dto.EmailAddress{
		Address: "from@bar.com",
		Name:    "From Bar",
	}

	suite.APIServiceMock.On("GetMessage", mock.Anything).Return(getMessageFixture, nil).Times(1)
	suite.APIServiceMock.On("UpdateMessage", mock.Anything, mock.Anything).Return(nil).Times(1)

	cmd := NewCmdGet(suite.CmdOptions)
	cmd.SetArgs([]string{"mess001"})
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Regexp(suite.T(), fmt.Sprintf(`Id[ ]*:[ ]*%s`, getMessageFixture.ID), cmdOutput)
	assert.Regexp(suite.T(), `From[ ]*:[ ]*from\@bar\.com \(From Bar\)`, cmdOutput)
	assert.Regexp(suite.T(), `Cc[ ]*:[ ]*cc\@bar\.com \(Cc Bar\)`, cmdOutput)
	assert.Regexp(suite.T(), `Bcc[ ]*:[ ]*bcc\@bar\.com \(Bcc Bar\)`, cmdOutput)
	assert.Regexp(suite.T(), fmt.Sprintf(`Subject[ ]*:[ ]*%s`, getMessageFixture.Subject), cmdOutput)
	assert.Regexp(
		suite.T(),
		fmt.Sprintf(`Retention Date[ ]*:[ ]*%s`, getMessageFixture.RetentionDate.Format("02 January 2006 15:04:05")),
		cmdOutput,
	)
	assert.Regexp(
		suite.T(),
		fmt.Sprintf(`Created At[ ]*:[ ]*%s`, getMessageFixture.CreatedAt.Format("02 January 2006 15:04:05")),
		cmdOutput,
	)
	assert.Contains(suite.T(), cmdOutput, getMessageFixture.Text)
}

func (suite *GetCmdSuite) TestNewCmdGet_WithAttachments() {
	getMessageFixture := &dto.Message{
		Attachments: []dto.MessageAttachment{
			{
				ID:       "ATTACH000001",
				Filename: "happy.png",
			},
			{
				ID:       "ATTACH000002",
				Filename: "cv.pdf",
			},
		},
	}
	getMessageFixture.HasAttachments = true

	suite.APIServiceMock.On("GetMessage", mock.Anything).Return(getMessageFixture, nil).Times(1)
	suite.APIServiceMock.On("UpdateMessage", mock.Anything, mock.Anything).Return(nil).Times(1)

	cmd := NewCmdGet(suite.CmdOptions)
	cmd.SetArgs([]string{"mess001"})
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Contains(suite.T(), cmdOutput, "Attachments:")
	assert.Contains(
		suite.T(),
		cmdOutput,
		fmt.Sprintf("• %s - %s", getMessageFixture.Attachments[0].ID, getMessageFixture.Attachments[0].Filename),
	)
	assert.Contains(
		suite.T(),
		cmdOutput,
		fmt.Sprintf("• %s - %s", getMessageFixture.Attachments[1].ID, getMessageFixture.Attachments[1].Filename),
	)
}

func (suite *GetCmdSuite) TestNewCmdGet_NotFound() {
	getMessageFixture := &dto.Message{}
	getMessageFixtureError := &errors.HTTPError{
		Code:   404,
		Detail: "Not Found",
	}

	suite.APIServiceMock.On("GetMessage", mock.Anything).Return(getMessageFixture, getMessageFixtureError).Times(1)

	cmd := NewCmdGet(suite.CmdOptions)
	cmd.SetArgs([]string{"mess002"})
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Contains(suite.T(), cmdOutput, "Error: Not Found [404]")
}

func (suite *GetCmdSuite) TestNewCmdGet_UnableToUpdateMessage() {
	getMessageFixture := &dto.Message{}
	updateMessageFixtureError := &errors.HTTPError{
		Code:   500,
		Detail: "Internal Server Error",
	}

	suite.APIServiceMock.On("GetMessage", mock.Anything).Return(getMessageFixture, nil).Times(1)
	suite.APIServiceMock.On("UpdateMessage", mock.Anything, mock.Anything).Return(updateMessageFixtureError).Times(1)

	cmd := NewCmdGet(suite.CmdOptions)
	cmd.SetArgs([]string{"mess002"})
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Contains(suite.T(), cmdOutput, "Error: Internal Server Error [500]")
}
