package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type InstallDlvCmd struct{}

func (c *InstallDlvCmd) Run(cmd *cobra.Command, args []string) error {
	installCmd := exec.CommandContext(
		cmd.Context(),
		"go",
		"install",
		"-x", "-v",
		"github.com/go-delve/delve/cmd/dlv@latest",
	)
	installCmd.Stderr = os.Stderr
	installCmd.Stdout = os.Stdout
	return installCmd.Run()
}
