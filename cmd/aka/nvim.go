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
	installCmd.Flags().BoolVar(&install.InstallNodeJSForCoC, "install-nodejs-for-coc", false, "install nodejs for coc")
	installCmd.Flags().StringVar(&install.NodeJSVersionForCoC, "nodejs-version-for-coc", "v22.14.0", "nodejs version for coc")

	cmd.AddCommand(installCmd)

	return cmd
}
