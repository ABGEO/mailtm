package account

import (
	"github.com/abgeo/mailtm/pkg/command"
	"github.com/spf13/cobra"
)

func NewCmd(options command.Options) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "account",
		Short: "Manage Account",
	}

	cmds.AddCommand(NewCmdCurrent(options))

	return cmds
}
