package main

import(
	"github.com/spf13/cobra"

	cmdpkg "github.com/pigfall/aka/pkg/cmd"
)

func dockerCmd()*cobra.Command{
  c := &cobra.Command{
    Use:"docker",
  }

  install := cmdpkg.InstallDockerCmd{}
  installCmd := &cobra.Command{
    Use:"install",
    RunE: install.Run,
  }

  c.AddCommand(
      installCmd,
  )

  return c

}
