package auth

import (
	"github.com/abgeo/mailtm/pkg/util"
	"github.com/spf13/cobra"
)

func NewCmd(options util.CmdOptions) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "auth",
		Short: "Manage Authentication",
	}

	cmds.AddCommand(NewCmdRandom(options))

	return cmds
}
