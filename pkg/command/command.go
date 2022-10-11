package command

import (
	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/abgeo/mailtm/pkg/types"
	"github.com/spf13/cobra"
)

type Options struct {
	Version    types.Version
	Config     configs.Config
	APIService service.APIServiceInterface
	SSEService *service.SSEService
}

type Command interface {
	Complete(cmd *cobra.Command, args []string) error
	Validate() error
	Run() error
}

func GetRunner(command Command) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(command.Complete(cmd, args))
		cobra.CheckErr(command.Validate())
		cobra.CheckErr(command.Run())
	}
}
