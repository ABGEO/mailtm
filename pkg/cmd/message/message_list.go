package message

import (
	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/dto"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/abgeo/mailtm/pkg/util"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type CommandList struct {
	Config  configs.Config
	Service *service.APIService
}

func NewCmdList(options util.CmdOptions) *cobra.Command {
	opts := &CommandList{
		Config:  options.Config,
		Service: service.NewAPIService(),
	}
	opts.Service.SetToken(&dto.Token{
		ID:    options.Config.Auth.ID,
		Token: options.Config.Auth.Token,
	})

	cmds := &cobra.Command{
		Use:   "list",
		Short: "List messages",
		Args:  cobra.NoArgs,
		Run:   util.GetCmdRunner(opts),
	}

	return cmds
}

func (command *CommandList) Complete(cmd *cobra.Command, args []string) error { return nil }

func (command *CommandList) Validate() error { return nil }

func (command *CommandList) Run() error {
	messages, err := command.Service.GetMessages()
	if err != nil {
		return err
	}

	tableData := pterm.TableData{
		{
			"ID",
			"Seen",
			"Subject",
			"From",
			"Intro",
		},
	}

	for _, message := range messages {
		seen := "No"
		if message.Seen {
			seen = "Yes"
		}

		tableData = append(tableData, []string{
			message.ID,
			seen,
			message.Subject,
			util.EmailAddressesToString(message.From),
			message.Intro,
		})
	}

	_ = pterm.DefaultTable.
		WithData(tableData).
		WithHeaderRowSeparator("-").
		WithHasHeader().
		WithLeftAlignment().
		Render()

	return nil
}
