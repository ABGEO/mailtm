package command

import (
	"io"

	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/errors"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/abgeo/mailtm/pkg/types"
	"github.com/spf13/cobra"
)

type Options struct {
	Writer     io.Writer
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
		errors.CheckErr(command.Complete(cmd, args), 1)
		errors.CheckErr(command.Validate(), 1)
		errors.CheckErr(command.Run(), 1)
	}
}
