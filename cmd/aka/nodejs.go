package main

import (
	"github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func nodejsCmd()*cobra.Command{
  cobraCmd := cobra.Command{
    Use:"nodejs",
  }

  installNodejs := &cmd.NodejsInstallCmd{

  }
  installNodejsCmd := cobra.Command{
    Use:"install",
    RunE: installNodejs.Run,
  }
  installNodejsCmd.Flags().StringVar(
      &installNodejs.Version,
      "version",
      "v22.14.0",
      "nodejs version",
  )

  cobraCmd.AddCommand(
      &installNodejsCmd,
  )

  return &cobraCmd
}
