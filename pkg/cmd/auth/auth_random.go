package auth

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/abgeo/mailtm/pkg/util"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	addressLength  = 10
	passwordLength = 8
)

type CommandRand struct {
	Writer  io.Writer
	Config  configs.Config
	Service service.APIServiceInterface
}

func NewCmdRandom(options command.Options) *cobra.Command {
	opts := &CommandRand{
		Writer:  options.Writer,
		Config:  options.Config,
		Service: options.APIService,
	}

	cmds := &cobra.Command{
		Use:   "random",
		Short: "Authenticate with random email",
		Args:  cobra.NoArgs,
		Run:   command.GetRunner(opts),
	}

	return cmds
}

func (command *CommandRand) Complete(cmd *cobra.Command, args []string) error { return nil }

func (command *CommandRand) Validate() error { return nil }

func (command *CommandRand) Run() error {
	account, password, err := command.createRandomAccount()
	if err != nil {
		return err
	}

	token, err := command.Service.GetToken(dto.Credentials{
		Address:  account.Address,
		Password: password,
	})
	if err != nil {
		return err
	}

	command.saveAuthData(token, account)

	return command.printAccountInfo(account, password)
}

func (command *CommandRand) createRandomAccount() (account *dto.Account, password string, err error) {
	domains, err := command.Service.GetDomains()
	if err != nil {
		return account, password, err
	}

	randomDomainIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(domains))))
	domain := domains[randomDomainIndex.Int64()]
	address := fmt.Sprintf("%s@%s", strings.ToLower(util.RandomString(addressLength)), domain.Domain)
	password = util.RandomString(passwordLength)

	account, err = command.Service.CreateAccount(dto.AccountWrite{
		Address:  address,
		Password: password,
	})
	if err != nil {
		return account, password, err
	}

	return account, password, nil
}

func (command *CommandRand) saveAuthData(token *dto.Token, account *dto.Account) {
	command.Config.Auth.ID = token.ID
	command.Config.Auth.Email = account.Address
	command.Config.Auth.Token = token.Token
	command.Config.Write()
}

func (command *CommandRand) printAccountInfo(account *dto.Account, password string) (err error) {
	pterm.DefaultBasicText.
		WithWriter(command.Writer).
		Println("New random account has been created and authenticated")

	return pterm.DefaultTable.
		WithWriter(command.Writer).
		WithData(pterm.TableData{
			{"Email", account.Address},
			{"Password", password},
		}).
		WithSeparator(" : ").
		WithBoxed().
		WithLeftAlignment().
		Render()
}
