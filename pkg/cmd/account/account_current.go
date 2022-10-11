package account

import (
	"fmt"
	"io"

	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type CommandCurrent struct {
	Writer  io.Writer
	Service service.APIServiceInterface
}

func NewCmdCurrent(options command.Options) *cobra.Command {
	opts := &CommandCurrent{
		Writer:  options.Writer,
		Service: options.APIService,
	}
	opts.Service.SetToken(&dto.Token{
		ID:    options.Config.Auth.ID,
		Token: options.Config.Auth.Token,
	})

	cmds := &cobra.Command{
		Use:   "current",
		Short: "Get current account",
		Args:  cobra.NoArgs,
		Run:   command.GetRunner(opts),
	}

	return cmds
}

func (command *CommandCurrent) Complete(cmd *cobra.Command, args []string) error { return nil }

func (command *CommandCurrent) Validate() error { return nil }

func (command *CommandCurrent) Run() error {
	data, err := command.getTableData()
	if err != nil {
		return err
	}

	return pterm.DefaultTable.
		WithWriter(command.Writer).
		WithData(data).
		WithSeparator(" : ").
		WithBoxed().
		WithLeftAlignment().
		Render()
}

func (command *CommandCurrent) getTableData() (data pterm.TableData, err error) {
	const mbShift = 20

	account, err := command.Service.GetCurrentAccount()
	if err != nil {
		return data, err
	}

	quota := float64(account.Quota) / (1 << mbShift)
	used := float64(account.Used) / (1 << mbShift)

	return pterm.TableData{
		{"ID", account.ID},
		{"Email", account.Address},
		{"Usage", fmt.Sprintf("%.2f MB / %.2f MB", used, quota)},
		{"Created At", account.CreatedAt.Local().Format("02 January 2006 15:04:05")},
	}, err
}
