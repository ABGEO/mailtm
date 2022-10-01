package cmd

import (
	"fmt"

	"github.com/abgeo/mailtm/pkg/cmd/account"
	"github.com/abgeo/mailtm/pkg/cmd/auth"
	"github.com/abgeo/mailtm/pkg/cmd/message"
	"github.com/abgeo/mailtm/pkg/util"
	"github.com/spf13/cobra"
)

const (
	VersionName = "Apple"
	Version     = "0.1.0"
)

func NewCmd(options util.CmdOptions) *cobra.Command {
	cmds := &cobra.Command{
		Use:     "mailtm",
		Version: fmt.Sprintf("%s (%s)", Version, VersionName),
		Short:   "CLI client for Mail.tm disposable mail service.",
		Long: `		   _  _    _              
 _ __ ___    __ _ (_)| |  | |_  _ __ ___  
| '_ ' _ \  / _' || || |  | __|| '_ ' _ \
| | | | | || (_| || || | _| |_ | | | | | |
|_| |_| |_| \__,_||_||_|(_)\__||_| |_| |_|

CLI client for Mail.tm disposable mail service.
`,
	}

	cmds.AddCommand(account.NewCmd(options))
	cmds.AddCommand(auth.NewCmd(options))
	cmds.AddCommand(message.NewCmd(options))

	return cmds
}
