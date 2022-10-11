package command

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var errComplete = errors.New("Complete.Error")

type dummyCommand struct{}

func (command *dummyCommand) Complete(cmd *cobra.Command, args []string) error {
	if args[0] == "fail" {
		return errComplete
	}

	return nil
}

func (command *dummyCommand) Validate() error { return nil }

func (command *dummyCommand) Run() error { return nil }

type CMDSuite struct {
	suite.Suite
}

func TestCMDSuite(t *testing.T) {
	suite.Run(t, new(CMDSuite))
}

func (suite *CMDSuite) TestGetRunner() {
	runner := GetRunner(&dummyCommand{})
	cmds := &cobra.Command{
		Use: "dummy",
		Run: runner,
	}

	runner(cmds, []string{"pass"})

	if os.Getenv("CRASH") == "1" {
		runner(cmds, []string{"fail"})
	}

	//nolint:gosec
	cmd := exec.Command(os.Args[0], "-test.run=TestCMDSuite", "-testify.m=TestGetRunner")
	cmd.Env = append(os.Environ(), "CRASH=1")
	err := cmd.Run()

	var e *exec.ExitError

	errors.As(err, &e)

	assert.Equal(suite.T(), 1, e.ExitCode())
	assert.Equal(suite.T(), false, e.Success())
}
