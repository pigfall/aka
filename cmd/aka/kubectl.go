package main

import (
	cmdpkg "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func kubectlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "kubectl",
	}

	install := cmdpkg.KubectlInstallCmd{}
	installCmd := &cobra.Command{
		Use:  "install",
		RunE: install.Run,
	}

	cmd.AddCommand(
		installCmd,
	)

	return cmd
}
