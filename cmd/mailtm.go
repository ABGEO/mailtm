package main

import (
	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/cmd"
	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/abgeo/mailtm/pkg/types"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "unknown" //nolint:gochecknoglobals
	date    = "unknown" //nolint:gochecknoglobals
)

func main() {
	appVersion := types.Version{
		Number: version,
		Commit: commit,
		Date:   date,
	}
	config := configs.NewConfig()
	cmdOpts := command.Options{
		Version:    appVersion,
		Config:     config,
		APIService: service.NewAPIService(appVersion),
		SSEService: service.NewSSEService(appVersion, config.Auth.AuthConfig),
	}
	rootCmd := cmd.NewCmd(cmdOpts)
	cobra.CheckErr(rootCmd.Execute())
}
