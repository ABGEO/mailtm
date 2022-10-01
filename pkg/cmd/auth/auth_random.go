package auth

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/abgeo/mailtm/configs"
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
	Config  configs.Config
	Service *service.APIService
}

func NewCmdRandom(options util.CmdOptions) *cobra.Command {
	opts := &CommandRand{
		Config:  options.Config,
		Service: service.NewAPIService(options.Version),
	}

	cmds := &cobra.Command{
		Use:   "random",
		Short: "Authenticate with random email",
		Args:  cobra.NoArgs,
		Run:   util.GetCmdRunner(opts),
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

	command.Config.Auth.ID = token.ID
	command.Config.Auth.Email = account.Address
	command.Config.Auth.Token = token.Token
	command.Config.Write()

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

func (command *CommandRand) printAccountInfo(account *dto.Account, password string) (err error) {
	pterm.DefaultBasicText.Println("New random account has been created and authenticated")

	return pterm.DefaultTable.
		WithData(pterm.TableData{
			{"Email", account.Address},
			{"Password", password},
		}).
		WithSeparator(" : ").
		WithBoxed().
		WithLeftAlignment().
		Render()
}
