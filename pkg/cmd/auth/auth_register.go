package auth

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"

	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/errors"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type CommandRegister struct {
	Writer  io.Writer
	Config  configs.Config
	Service service.APIServiceInterface

	Username string
	Password string

	Domain string
	Login  bool

	domains []string
}

func NewCmdRegister(options command.Options) *cobra.Command {
	const numberOfArguments = 2

	opts := &CommandRegister{
		Writer:  options.Writer,
		Config:  options.Config,
		Service: options.APIService,
	}

	cmds := &cobra.Command{
		Use:   "register [username] [password]",
		Short: "Register new account",
		Example: `mailtm auth register john Pa$$w0rd
mailtm auth register john Pa$$w0rd --domain doe.com`,
		Args: cobra.ExactArgs(numberOfArguments),
		Run:  command.GetRunner(opts),
	}

	cmds.Flags().StringVarP(
		&opts.Domain,
		"domain",
		"d",
		"",
		"The domain to generate the email address. If not specified, random one will be used.",
	)
	cmds.Flags().BoolVarP(
		&opts.Login,
		"login",
		"l",
		false,
		"Whether login after creating an account or not",
	)

	return cmds
}

func (command *CommandRegister) Complete(cmd *cobra.Command, args []string) error {
	command.Username = args[0]
	command.Password = args[1]

	domains, err := command.Service.GetDomains()
	if err != nil {
		return err
	}

	for _, domain := range domains {
		command.domains = append(command.domains, domain.Domain)
	}

	if command.Domain == "" {
		randomDomainIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(command.domains))))
		command.Domain = command.domains[randomDomainIndex.Int64()]
	}

	return nil
}

func (command *CommandRegister) Validate() error {
	for _, domain := range command.domains {
		if domain == command.Domain {
			return nil
		}
	}

	return errors.NewInvalidDomainError(command.domains, command.Domain)
}

func (command *CommandRegister) Run() error {
	account, err := command.Service.CreateAccount(dto.AccountWrite{
		Address:  fmt.Sprintf("%s@%s", command.Username, command.Domain),
		Password: command.Password,
	})
	if err != nil {
		return err
	}

	pterm.DefaultBasicText.
		WithWriter(command.Writer).
		Printfln("New account %s has been created successfully.", account.Address)

	if command.Login {
		if err = command.logIn(account, command.Password); err != nil {
			return err
		}

		pterm.DefaultBasicText.
			WithWriter(command.Writer).
			Printfln("Authenticated with new account %s.", account.Address)
	}

	return nil
}

func (command *CommandRegister) logIn(account *dto.Account, password string) error {
	token, err := command.Service.GetToken(dto.Credentials{
		Address:  account.Address,
		Password: password,
	})
	if err != nil {
		return err
	}

	command.Config.Auth.ID = token.ID
	command.Config.Auth.Email = account.Address
	command.Config.Auth.Token = token.Token
	command.Config.Write()

	return nil
}
