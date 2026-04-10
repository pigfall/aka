package main

import (
	cmdpkg "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func opencodeCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "opencode",
	}

	install := cmdpkg.InstallOpenCodeCmd{}
	installCmd := &cobra.Command{
		Use:  "install",
		RunE: install.Run,
	}
	installCmd.Flags().StringVar(
		&install.Version,
		"version",
		"",
		"version to install (default: latest bundled version)",
	)

	c.AddCommand(
		installCmd,
	)

	return c
}
