package main

import (
	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/cmd"
	"github.com/abgeo/mailtm/pkg/util"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "unknown" //nolint:gochecknoglobals
	date    = "unknown" //nolint:gochecknoglobals
)

func main() {
	cmdOpts := util.CmdOptions{
		Version: util.Version{
			Number: version,
			Commit: commit,
			Date:   date,
		},
		Config: configs.NewConfig(),
	}
	rootCmd := cmd.NewCmd(cmdOpts)
	cobra.CheckErr(rootCmd.Execute())
}
