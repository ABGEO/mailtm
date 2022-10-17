package account

import (
	"io"

	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type CommandDelete struct {
	Writer  io.Writer
	Service service.APIServiceInterface
	Config  configs.Config
}

func NewCmdDelete(options command.Options) *cobra.Command {
	opts := &CommandDelete{
		Writer:  options.Writer,
		Service: options.APIService,
		Config:  options.Config,
	}
	opts.Service.SetToken(&dto.Token{
		ID:    options.Config.Auth.ID,
		Token: options.Config.Auth.Token,
	})

	cmds := &cobra.Command{
		Use:   "delete",
		Short: "Delete the current account",
		Args:  cobra.NoArgs,
		Run:   command.GetRunner(opts),
	}

	return cmds
}

func (command *CommandDelete) Complete(cmd *cobra.Command, args []string) error { return nil }

func (command *CommandDelete) Validate() error { return nil }

func (command *CommandDelete) Run() error {
	confirm, _ := pterm.DefaultInteractiveConfirm.
		WithDefaultText("Are you sure you want to delete the account?").
		Show()

	if confirm {
		err := command.Service.RemoveAccount(command.Config.Auth.ID)
		if err != nil {
			return err
		}

		command.Config.Auth.ID = ""
		command.Config.Auth.Email = ""
		command.Config.Auth.Token = ""
		command.Config.Write()

		pterm.DefaultParagraph.
			WithWriter(command.Writer).
			Println("Account has been deleted successfully.")
	}

	return nil
}
