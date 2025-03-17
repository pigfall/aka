package main

import (
	pkgcmd "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func reactCmd() *cobra.Command {
	cobraCmd := cobra.Command{
		Use: "react",
	}

	initReactUILibrary := pkgcmd.InitReactUILibraryCmd{}
	initReactUILibraryCmd := &cobra.Command{
		Use:  "init-ui-library",
		RunE: initReactUILibrary.Run,
	}

	cobraCmd.AddCommand(
		initReactUILibraryCmd,
	)

	return &cobraCmd
}
