package message

import (
	"github.com/abgeo/mailtm/pkg/command"
	"github.com/spf13/cobra"
)

func NewCmd(options command.Options) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "message",
		Short: "Manage Messages",
	}

	cmds.AddCommand(NewCmdDownloadAttachment(options))
	cmds.AddCommand(NewCmdGet(options))
	cmds.AddCommand(NewCmdList(options))

	return cmds
}
