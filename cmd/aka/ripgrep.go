package main

import (
	pkgcmd "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func ripgrepCmd() *cobra.Command {
	cobraCmd := cobra.Command{
		Use: "ripgrep",
	}

	install := pkgcmd.RipgrepInstallCmd{}
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install ripgrep",
		RunE:  install.Run,
	}

	cobraCmd.AddCommand(
		installCmd,
	)

	return &cobraCmd
}
