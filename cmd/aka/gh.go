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

	c.AddCommand(
		installCmd,
	)

	return c

}
