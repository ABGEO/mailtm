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
	Config     configs.Config
	Service    *service.APIService
	SSOService *service.SSEService

	Watch bool
}

func NewCmdList(options util.CmdOptions) *cobra.Command {
	opts := &CommandList{
		Config:     options.Config,
		Service:    service.NewAPIService(options.Version),
		SSOService: service.NewSSEService(options.Version, options.Config.Auth.AuthConfig),
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

	cmds.Flags().BoolVarP(&opts.Watch, "watch", "w", false, "Watch new messages")

	return cmds
}

func (command *CommandList) Complete(cmd *cobra.Command, args []string) error { return nil }

func (command *CommandList) Validate() error { return nil }

func (command *CommandList) Run() error {
	tableHeader := []string{
		"ID",
		"Seen",
		"Subject",
		"From",
		"Intro",
	}

	tableData, err := command.getInitialTableRows()
	if err != nil {
		return err
	}

	area, _ := pterm.DefaultArea.Start()
	command.drawTableInArea(area, tableHeader, tableData)

	if command.Watch {
		command.watchMessages(area, tableHeader, tableData)
	}

	return area.Stop()
}

func (command *CommandList) getInitialTableRows() (rows pterm.TableData, err error) {
	messages, err := command.Service.GetMessages()
	if err != nil {
		return rows, err
	}

	for _, message := range messages {
		rows = append(rows, command.messageToTableRow(message))
	}

	return rows, nil
}

func (command *CommandList) messageToTableRow(message dto.MessagesItem) []string {
	seen := "No"
	if message.Seen {
		seen = "Yes"
	}

	return []string{
		message.ID,
		seen,
		message.Subject,
		util.EmailAddressesToString(message.From),
		message.Intro,
	}
}

func (command *CommandList) drawTableInArea(area *pterm.AreaPrinter, tableHeader []string, tableData pterm.TableData) {
	tableData = append(pterm.TableData{tableHeader}, tableData...)

	table, _ := pterm.DefaultTable.
		WithData(tableData).
		WithHeaderRowSeparator("-").
		WithHasHeader().
		WithLeftAlignment().
		Srender()

	area.Update(table)
}

func (command *CommandList) watchMessages(area *pterm.AreaPrinter, tableHeader []string, tableData pterm.TableData) {
	// @todo: We also receive an event when message is seen. We have to fix it.
	_ = command.SSOService.SubscribeMessages(command.Config.Auth.ID, func(message dto.MessagesItem) {
		tableData = append(pterm.TableData{command.messageToTableRow(message)}, tableData...)

		command.drawTableInArea(area, tableHeader, tableData)
	})
}
