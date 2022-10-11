package test

import (
	"bytes"
	"io"

	"github.com/abgeo/mailtm/pkg/command"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type BaseCMDSuite struct {
	suite.Suite

	CmdOptions command.Options
}

func (suite *BaseCMDSuite) GetCommandOutput(command *cobra.Command) string {
	outBuffer := bytes.NewBufferString("")

	command.SetOut(outBuffer)

	if err := command.Execute(); err != nil {
		return ""
	}

	out, err := io.ReadAll(outBuffer)
	if err != nil {
		return ""
	}

	return string(out)
}
