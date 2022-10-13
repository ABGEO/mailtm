package command

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var errComplete = errors.New("Complete.Error")

type dummyCommand struct{}

func (command *dummyCommand) Complete(cmd *cobra.Command, args []string) error {
	return errComplete
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

	fakeExit := func(code int) { panic(fmt.Sprintf("exited with %d", code)) }
	patch := monkey.Patch(os.Exit, fakeExit)

	defer patch.Unpatch()

	assert.PanicsWithValue(suite.T(), "exited with 1", func() { runner(cmds, []string{}) })
}
