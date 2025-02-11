package main

import (
	cmdpkg "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func nvimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "nvim",
	}

	install := cmdpkg.NvimInstallCmd{}
	installCmd := &cobra.Command{
		Use:  "install",
		RunE: install.Run,
	}
	installCmd.Flags().BoolVar(&install.InstallPlugin, "plugin", true, "install plugin")

	cmd.AddCommand(installCmd)

	return cmd
}
