package main

import (
	cmdpkg "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func nvmCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "nvm",
	}

	installNvm := cmdpkg.NvmInstallCmd{}
	installNvmCmd := &cobra.Command{
		Use:   "install",
		Short: "Install nvm",
		Args:  cobra.NoArgs,
		RunE:  installNvm.Run,
	}

	nodejsCmd := &cobra.Command{
		Use: "nodejs",
	}

	installNodejs := cmdpkg.NvmNodejsInstallCmd{}
	installNodejsCmd := &cobra.Command{
		Use:   "install <nodejs-version>",
		Short: "Install nodejs by nvm",
		Args:  cobra.ExactArgs(1),
		RunE:  installNodejs.Run,
	}
	installNodejsCmd.Example = "aka nvm nodejs install v22.14.0"

	list := cmdpkg.NvmListCmd{}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List nodejs versions by nvm",
		Args:  cobra.NoArgs,
		RunE:  list.Run,
	}

	defaultVersion := cmdpkg.NvmNodejsDefaultCmd{}
	defaultVersionCmd := &cobra.Command{
		Use:   "default [nodejs-version]",
		Short: "Set default nodejs version by nvm",
		Args:  cobra.MaximumNArgs(1),
		RunE:  defaultVersion.Run,
	}
	defaultVersionCmd.Example = "aka nvm nodejs default v22.14.0"

	nodejsCmd.AddCommand(
		installNodejsCmd,
		listCmd,
		defaultVersionCmd,
	)

	c.AddCommand(
		installNvmCmd,
		nodejsCmd,
	)

	return c
}
