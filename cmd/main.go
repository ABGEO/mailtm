package main

import (
	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/cmd"
	"github.com/abgeo/mailtm/pkg/util"
	"github.com/spf13/cobra"
)

func main() {
	cmdOpts := util.CmdOptions{
		Config: configs.NewConfig(),
	}
	rootCmd := cmd.NewCmd(cmdOpts)
	cobra.CheckErr(rootCmd.Execute())
}
