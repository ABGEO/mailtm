package account

import (
	"github.com/abgeo/mailtm/pkg/util"
	"github.com/spf13/cobra"
)

func NewCmd(options util.CmdOptions) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "account",
		Short: "Manage Account",
	}

	cmds.AddCommand(NewCmdCurrent(options))

	return cmds
}
