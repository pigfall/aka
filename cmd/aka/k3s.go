package main

import (
	cmdpkg "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func k3sCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "k3s",
	}

	k3sInstall := &cmdpkg.K3SInstallCmd{}
	k3sInstallCmd := cobra.Command{
		Use:  "install",
		RunE: k3sInstall.Run,
	}

	cmd.AddCommand(
		&k3sInstallCmd,
	)

	return cmd
}
