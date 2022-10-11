package test

import (
	"bytes"
	"regexp"

	"github.com/abgeo/mailtm/pkg/command"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type BaseCMDSuite struct {
	suite.Suite

	Buffer     *bytes.Buffer
	CmdOptions command.Options
}

func (suite *BaseCMDSuite) GetCommandOutput(command *cobra.Command) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)" +
		"?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

	outBuffer := suite.Buffer
	if outBuffer == nil {
		outBuffer = bytes.NewBufferString("")
	}

	command.SetOut(outBuffer)

	if err := command.Execute(); err != nil {
		return "<err>"
	}

	return regexp.MustCompile(ansi).ReplaceAllString(outBuffer.String(), "")
}
