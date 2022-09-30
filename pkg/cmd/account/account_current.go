package account

import (
	"fmt"

	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/abgeo/mailtm/pkg/util"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type CommandCurrent struct {
	Service *service.APIService
}

func NewCmdCurrent(options util.CmdOptions) *cobra.Command {
	opts := &CommandCurrent{
		Service: service.NewAPIService(),
	}
	opts.Service.SetToken(&dto.Token{
		ID:    options.Config.Auth.ID,
		Token: options.Config.Auth.Token,
	})

	cmds := &cobra.Command{
		Use:   "current",
		Short: "Get current account",
		Args:  cobra.NoArgs,
		Run:   util.GetCmdRunner(opts),
	}

	return cmds
}

func (command *CommandCurrent) Complete(cmd *cobra.Command, args []string) error { return nil }

func (command *CommandCurrent) Validate() error { return nil }

func (command *CommandCurrent) Run() error {
	const mbShift = 20

	account, err := command.Service.GetCurrentAccount()
	if err != nil {
		return err
	}

	quota := float64(account.Quota) / (1 << mbShift)
	used := float64(account.Used) / (1 << mbShift)

	return pterm.DefaultTable.
		WithData(pterm.TableData{
			{"ID", account.ID},
			{"Email", account.Address},
			{"Usage", fmt.Sprintf("%.2f MB / %.2f MB", used, quota)},
			{"Created At", account.CreatedAt.Local().Format("02 January 2006 15:04:05")},
		}).
		WithSeparator(" : ").
		WithBoxed().
		WithLeftAlignment().
		Render()
}
