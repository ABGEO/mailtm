package main

import (
	"os"

	"github.com/abgeo/mailtm/configs"
	"github.com/abgeo/mailtm/pkg/cmd"
	"github.com/abgeo/mailtm/pkg/command"
	"github.com/abgeo/mailtm/pkg/errors"
	"github.com/abgeo/mailtm/pkg/service"
	"github.com/abgeo/mailtm/pkg/types"
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
		Writer:     os.Stdout,
		Version:    appVersion,
		Config:     config,
		APIService: service.NewAPIService(appVersion),
		SSEService: service.NewSSEService(appVersion, config.Auth.AuthConfig),
	}
	rootCmd := cmd.NewCmd(cmdOpts)
	errors.CheckErr(rootCmd.Execute(), 1)
}
