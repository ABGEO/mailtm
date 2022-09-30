package util

import (
	"github.com/abgeo/mailtm/configs"
	"github.com/spf13/cobra"
)

type CmdOptions struct {
	Config configs.Config
}

type Command interface {
	Complete(cmd *cobra.Command, args []string) error
	Validate() error
	Run() error
}

func GetCmdRunner(command Command) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(command.Complete(cmd, args))
		cobra.CheckErr(command.Validate())
		cobra.CheckErr(command.Run())
	}
}
