package main

import (
	pkgcmd "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func piCmd() *cobra.Command {
	cobraCmd := cobra.Command{
		Use: "pi",
	}

	install := pkgcmd.PiInstallCmd{}
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install the Pi Coding Agent CLI",
		RunE:  install.Run,
	}
	installCmd.Flags().StringVar(
		&install.Version,
		"version",
		"",
		"version to install, e.g. 1.6.0 (default: latest release)",
	)
	installCmd.Flags().BoolVar(
		&install.Force,
		"force",
		false,
		"reinstall even if the same pi version is already installed",
	)

	cobraCmd.AddCommand(
		installCmd,
	)

	return &cobraCmd
}
