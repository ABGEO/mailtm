package auth

import (
	"github.com/abgeo/mailtm/pkg/command"
	"github.com/spf13/cobra"
)

func NewCmd(options command.Options) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "auth",
		Short: "Manage Authentication",
	}

	cmds.AddCommand(NewCmdRandom(options))

	return cmds
}
