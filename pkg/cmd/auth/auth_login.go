package auth

import (
	"io"

	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const numberOfArguments = 2

type CommandLogin struct {
	Writer  io.Writer
	Config  configs.Config
	Service service.APIServiceInterface

	Email    string
	Password string
}

func NewCmdLogin(options command.Options) *cobra.Command {
	opts := &CommandLogin{
		Writer:  options.Writer,
		Config:  options.Config,
		Service: options.APIService,
	}

	cmds := &cobra.Command{
		Use:     "login [email] [password]",
		Short:   "Login with credentials",
		Example: "mailtm auth login john@doe.com Pa$$w0rd",
		Args:    cobra.ExactArgs(numberOfArguments),
		Run:     command.GetRunner(opts),
	}

	return cmds
}

func (command *CommandLogin) Complete(cmd *cobra.Command, args []string) error {
	command.Email = args[0]
	command.Password = args[1]

	return nil
}

func (command *CommandLogin) Validate() error { return nil }

func (command *CommandLogin) Run() error {
	token, err := command.Service.GetToken(dto.Credentials{
		Address:  command.Email,
		Password: command.Password,
	})
	if err != nil {
		return err
	}

	command.Config.Auth.ID = token.ID
	command.Config.Auth.Email = command.Email
	command.Config.Auth.Token = token.Token
	command.Config.Write()

	pterm.DefaultParagraph.
		WithWriter(command.Writer).
		Printfln("User %s has been authenticated successfully.", command.Email)

	return nil
}
