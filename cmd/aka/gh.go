package main

import (
	"github.com/spf13/cobra"

	cmdpkg "github.com/pigfall/aka/pkg/cmd"
)

func ghCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "gh",
	}

	install := cmdpkg.InstallGHCmd{}
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
