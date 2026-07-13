package main

import (
	cmdpkg "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func nvmCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "nvm",
	}

	install := cmdpkg.NvmInstallCmd{}
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install nvm",
		RunE:  install.Run,
	}
	installCmd.Flags().StringVar(
		&install.Version,
		"version",
		"v0.40.3",
		"nvm version",
	)
	installCmd.Flags().BoolVar(
		&install.Force,
		"force",
		false,
		"force reinstall nvm",
	)

	c.AddCommand(
		installCmd,
	)

	return c
}
