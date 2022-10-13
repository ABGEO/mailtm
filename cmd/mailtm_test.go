package main

import (
	"fmt"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MailtmSuite struct {
	suite.Suite
}

func TestMailtmSuite(t *testing.T) {
	suite.Run(t, new(MailtmSuite))
}

func (suite *MailtmSuite) TestExitCode() {
	fakeExit := func(code int) { panic(fmt.Sprintf("exited with %d", code)) }
	patch := monkey.Patch(os.Exit, fakeExit)

	defer patch.Unpatch()

	originalOsArgs := os.Args
	originalStdout := os.Stdout

	defer func(args []string, stdout *os.File) {
		os.Args = args
		os.Stdout = stdout
	}(originalOsArgs, originalStdout)

	os.Args = []string{"test", "wrong-command"}
	os.Stdout = os.NewFile(0, os.DevNull)

	assert.PanicsWithValue(suite.T(), "exited with 1", func() { main() })
}
