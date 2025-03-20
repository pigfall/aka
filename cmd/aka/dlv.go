package main

import (
	cmdpkg "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func dlvCmd() *cobra.Command {
	cobraCmd := cobra.Command{
		Use: "dlv",
	}

	install := &cmdpkg.InstallDlvCmd{}
	installCmd := &cobra.Command{
		Use:  "install",
		RunE: install.Run,
	}

	cobraCmd.AddCommand(
		installCmd,
	)
	return &cobraCmd
}
