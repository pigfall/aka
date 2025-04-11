package main

import (
	pkgcmd "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func tlstunCmd() *cobra.Command {
	cobraCmd := &cobra.Command{
		Use: "tlstun",
	}

	install := pkgcmd.TLSTunInstallCmd{}
	installCmd := &cobra.Command{
		Use:  "install",
		RunE: install.Run,
	}
	installCmd.Flags().StringVar(
		&install.Password,
		"password",
		"",
		"password",
	)

	clientRun := pkgcmd.TLSTunClientCmd{}
	clientRunCmd := &cobra.Command{
		Use:  "client",
		RunE: clientRun.Run,
	}
	clientRunCmd.Flags().StringVar(
		&clientRun.TLSTunPath,
		"tlstun-path",
		"tlstun",
		"tlstun path",
	)

	cobraCmd.AddCommand(
		installCmd,
		clientRunCmd,
	)

	return cobraCmd
}
