package message

import (
	"errors"
	"fmt"

	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/spf13/cobra"
)

type CommandDownloadAttachment struct {
	Config  configs.Config
	Service service.APIServiceInterface

	MessageID    string
	AttachmentID string

	FileDir string
}

var errAttachmentNotFound = errors.New("attachment not found")

func NewCmdDownloadAttachment(options command.Options) *cobra.Command {
	const numberOfArguments = 2

	opts := &CommandDownloadAttachment{
		Config:  options.Config,
		Service: options.APIService,
	}
	opts.Service.SetToken(&dto.Token{
		ID:    options.Config.Auth.ID,
		Token: options.Config.Auth.Token,
	})

	cmds := &cobra.Command{
		Use:     "download-attachment [message-id] [attachment-id]",
		Short:   "Download message attachment by ID",
		Example: "mailtm message download-attachment 63387xb206061bf4aaba863e ATTACH000001",
		Args:    cobra.ExactArgs(numberOfArguments),
		Run:     command.GetRunner(opts),
	}

	cmds.Flags().StringVarP(&opts.FileDir, "dir", "d", "./", "Directory path to save file in")

	return cmds
}

func (command *CommandDownloadAttachment) Complete(cmd *cobra.Command, args []string) error {
	command.MessageID = args[0]
	command.AttachmentID = args[1]

	return nil
}

func (command *CommandDownloadAttachment) Validate() error { return nil }

func (command *CommandDownloadAttachment) Run() error {
	message, err := command.Service.GetMessage(command.MessageID)
	if err != nil {
		return err
	}

	var attachment dto.MessageAttachment

	for _, item := range message.Attachments {
		if item.ID == command.AttachmentID {
			attachment = item

			break
		}
	}

	if attachment.ID == "" {
		return errAttachmentNotFound
	}

	return command.Service.DownloadMessageAttachment(
		command.MessageID,
		command.AttachmentID,
		fmt.Sprintf("%s/%s", command.FileDir, attachment.Filename),
	)
}
