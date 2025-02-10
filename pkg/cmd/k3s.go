package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type K3SInstallCmd struct {
}

func (c *K3SInstallCmd) Run(cobraCmd *cobra.Command, args []string) error {
	cmd := exec.Command(
		"bash",
		"-c",
		"curl -sfL https://get.k3s.io | sh -",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
