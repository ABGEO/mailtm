package message

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/errors"
	"github.com/abgeo/mailtm/test"
	"github.com/abgeo/mailtm/test/mocks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ListCmdSuite struct {
	test.BaseCMDSuite

	APIServiceMock *mocks.APIServiceInterface
	SSEServiceMock *mocks.SSEServiceInterface
}

func (suite *ListCmdSuite) SetupSuite() {
	suite.Buffer = bytes.NewBufferString("")
	suite.APIServiceMock = mocks.NewAPIServiceInterface(suite.T())
	suite.SSEServiceMock = mocks.NewSSEServiceInterface(suite.T())
	suite.CmdOptions = command.Options{
		Writer:     suite.Buffer,
		APIService: suite.APIServiceMock,
		SSEService: suite.SSEServiceMock,
	}

	errors.SetErrorHandler(func(msg interface{}, code int) {
		fmt.Fprintln(suite.Buffer, "Error:", msg)
	})

	suite.APIServiceMock.On("SetToken", mock.Anything)
}

func (suite *ListCmdSuite) TearDownTest() {
	suite.Buffer.Reset()
}

func TestListCmdSuite(t *testing.T) {
	suite.Run(t, new(ListCmdSuite))
}

func (suite *ListCmdSuite) TestNewCmdList_Success() {
	getMessagesFixture := dto.Messages{
		{Intro: "Message 1 Intro"},
		{Intro: "Message 2 Intro"},
	}
	getMessagesFixture[0].ID = "mess001"
	getMessagesFixture[0].Seen = false
	getMessagesFixture[0].Subject = "Message 1 Subject"
	getMessagesFixture[0].From = dto.EmailAddress{
		Address: "foo@bar.baz",
		Name:    "Foo Bar",
	}

	getMessagesFixture[1].ID = "mess002"
	getMessagesFixture[1].Seen = true
	getMessagesFixture[1].Subject = "Message 2 Subject"
	getMessagesFixture[1].From = dto.EmailAddress{
		Address: "bar@foo.baz",
		Name:    "Bar Foo",
	}

	suite.APIServiceMock.On("GetMessages", mock.Anything).Return(getMessagesFixture, nil).Times(1)

	cmd := NewCmdList(suite.CmdOptions)
	cmdOutput := suite.getCommandStdout(cmd)

	assert.Regexp(
		suite.T(),
		`mess001[ ]*|[ ]*No[ ]*|[ ]*Message 1 Subject[ ]*|[ ]*foo\@bar\.baz \(From Bar\)[ ]*|[ ]*Message 1 Intro`,
		cmdOutput,
	)
	assert.Regexp(
		suite.T(),
		`mess002[ ]*|[ ]*Yes[ ]*|[ ]*Message 2 Subject[ ]*|[ ]*bar\@foo\.baz \(Bar Foo\)[ ]*|[ ]*Message 2 Intro`,
		cmdOutput,
	)
}

func (suite *ListCmdSuite) TestNewCmdList_AccessDenied() {
	getMessagesFixture := dto.Messages{}
	getMessagesFixtureError := &errors.HTTPError{
		Code:   401,
		Detail: "Access Denied",
	}

	suite.APIServiceMock.On("GetMessages", mock.Anything).Return(getMessagesFixture, getMessagesFixtureError).Times(1)

	cmd := NewCmdList(suite.CmdOptions)
	cmdOutput := suite.GetCommandOutput(cmd)

	assert.Contains(suite.T(), cmdOutput, "Error: Access Denied [401]")
}

func (suite *ListCmdSuite) TestNewCmdList_WatchMode() {
	getMessagesFixture := dto.Messages{
		{Intro: "Message 1 Intro"},
		{Intro: "Message 2 Intro"},
	}
	getMessagesFixture[0].ID = "mess003"
	getMessagesFixture[0].Seen = false
	getMessagesFixture[0].Subject = "Message 1 Subject"
	getMessagesFixture[0].From = dto.EmailAddress{
		Address: "foo@bar.baz",
		Name:    "Foo Bar",
	}

	getMessagesFixture[1].ID = "mess004"
	getMessagesFixture[1].Seen = true
	getMessagesFixture[1].Subject = "Message 2 Subject"
	getMessagesFixture[1].From = dto.EmailAddress{
		Address: "bar@foo.baz",
		Name:    "Bar Foo",
	}

	suite.APIServiceMock.On("GetMessages", mock.Anything).Return(getMessagesFixture, nil).Times(1)
	suite.SSEServiceMock.On("SubscribeMessages", mock.Anything, mock.Anything).Return(nil).Times(1)

	cmd := NewCmdList(suite.CmdOptions)
	cmd.SetArgs([]string{"--watch"})
	cmdOutput := suite.getCommandStdout(cmd)

	assert.Regexp(
		suite.T(),
		`mess005[ ]*|[ ]*No[ ]*|[ ]*Message 3 Subject[ ]*|[ ]*baz\@foo\.bar \(Baz Foo\)[ ]*|[ ]*New Message`,
		cmdOutput,
	)
}

func (suite *ListCmdSuite) getCommandStdout(command *cobra.Command) string {
	originalStdout := os.Stdout
	readStream, writeStream, _ := os.Pipe()
	os.Stdout = writeStream

	_ = command.Execute()

	outChain := make(chan string)

	go func() {
		_, _ = io.Copy(suite.Buffer, readStream)
		outChain <- suite.Buffer.String()
	}()

	_ = writeStream.Close()

	os.Stdout = originalStdout

	return <-outChain
}
