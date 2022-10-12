package message

import (
	"fmt"
	"io"

	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/abgeo/mailtm/pkg/util"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type CommandGet struct {
	Writer  io.Writer
	Config  configs.Config
	Service service.APIServiceInterface

	ID string
}

func NewCmdGet(options command.Options) *cobra.Command {
	opts := &CommandGet{
		Writer:  options.Writer,
		Config:  options.Config,
		Service: options.APIService,
	}
	opts.Service.SetToken(&dto.Token{
		ID:    options.Config.Auth.ID,
		Token: options.Config.Auth.Token,
	})

	cmds := &cobra.Command{
		Use:   "get [id]",
		Short: "Get single message by ID",
		Args:  cobra.ExactArgs(1),
		Run:   command.GetRunner(opts),
	}

	return cmds
}

func (command *CommandGet) Complete(cmd *cobra.Command, args []string) error {
	command.ID = args[0]

	return nil
}

func (command *CommandGet) Validate() error { return nil }

func (command *CommandGet) Run() error {
	message, err := command.Service.GetMessage(command.ID)
	if err != nil {
		return err
	}

	err = command.Service.UpdateMessage(command.ID, dto.MessageWrite{
		Seen: true,
	})
	if err != nil {
		return err
	}

	command.printMessage(message)

	if message.HasAttachments {
		command.printAttachments(message)
	}

	return nil
}

func (command *CommandGet) printMessage(message *dto.Message) {
	_ = pterm.DefaultTable.
		WithWriter(command.Writer).
		WithData(command.messageToTableData(message)).
		WithSeparator(" : ").
		WithBoxed().
		WithLeftAlignment().
		Render()

	pterm.DefaultParagraph.WithWriter(command.Writer).Println()
	pterm.DefaultParagraph.WithWriter(command.Writer).Println(message.Text)
}

func (command *CommandGet) messageToTableData(message *dto.Message) (data pterm.TableData) {
	return pterm.TableData{
		{"Id", message.ID},
		{"From", util.EmailAddressesToString(message.From)},
		{"Cc", util.EmailAddressesToString(message.Cc...)},
		{"Bcc", util.EmailAddressesToString(message.Bcc...)},
		{"Subject", message.Subject},
		{"Retention Date", message.RetentionDate.Local().Format("02 January 2006 15:04:05")},
		{"Created At", message.CreatedAt.Local().Format("02 January 2006 15:04:05")},
	}
}

func (command *CommandGet) printAttachments(message *dto.Message) {
	pterm.DefaultParagraph.WithWriter(command.Writer).Println()
	pterm.DefaultParagraph.WithWriter(command.Writer).Println("Attachments:")
	pterm.DefaultParagraph.WithWriter(command.Writer).Println()

	attachments := make([]pterm.BulletListItem, len(message.Attachments))
	for i, attachment := range message.Attachments {
		attachments[i].Text = fmt.Sprintf("%s - %s", attachment.ID, attachment.Filename)
	}

	_ = pterm.DefaultBulletList.
		WithWriter(command.Writer).
		WithItems(attachments).
		Render()
}
