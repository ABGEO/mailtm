package message

import (
	"fmt"

	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/abgeo/mailtm/pkg/util"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type CommandGet struct {
	Config  configs.Config
	Service *service.APIService

	ID string
}

func NewCmdGet(options util.CmdOptions) *cobra.Command {
	opts := &CommandGet{
		Config:  options.Config,
		Service: service.NewAPIService(),
	}
	opts.Service.SetToken(&dto.Token{
		ID:    options.Config.Auth.ID,
		Token: options.Config.Auth.Token,
	})

	cmds := &cobra.Command{
		Use:   "get [id]",
		Short: "Get single message by ID",
		Args:  cobra.ExactArgs(1),
		Run:   util.GetCmdRunner(opts),
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

	err = pterm.DefaultTable.
		WithData(pterm.TableData{
			{"Id", message.ID},
			{"From", util.EmailAddressesToString(message.From)},
			{"Cc", util.EmailAddressesToString(message.Cc...)},
			{"Bcc", util.EmailAddressesToString(message.Bcc...)},
			{"Subject", message.Subject},
			{"Retention Date", message.RetentionDate.Local().Format("02 January 2006 15:04:05")},
			{"Created At", message.CreatedAt.Local().Format("02 January 2006 15:04:05")},
		}).
		WithSeparator(" : ").
		WithBoxed().
		WithLeftAlignment().
		Render()
	if err != nil {
		return err
	}

	pterm.DefaultParagraph.Println()
	pterm.DefaultParagraph.Println(message.Text)

	if message.HasAttachments {
		pterm.DefaultParagraph.Println()
		pterm.DefaultParagraph.Println("Attachments:")
		pterm.DefaultParagraph.Println()

		var attachments []pterm.BulletListItem
		for _, attachment := range message.Attachments {
			attachments = append(attachments, pterm.BulletListItem{
				Text: fmt.Sprintf("%s - %s", attachment.ID, attachment.Filename),
			})
		}

		if err = pterm.DefaultBulletList.WithItems(attachments).Render(); err != nil {
			return err
		}
	}

	return nil
}
