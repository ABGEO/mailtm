package message

import (
	"github.com/abgeo/mailtm/pkg/util"
	"github.com/spf13/cobra"
)

func NewCmd(options util.CmdOptions) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "message",
		Short: "Manage Messages",
	}

	cmds.AddCommand(NewCmdList(options))

	return cmds
}
