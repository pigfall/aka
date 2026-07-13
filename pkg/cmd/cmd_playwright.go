package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

type PlaywrightInitCmd struct {
}

func (c *PlaywrightInitCmd) Run(cmd *cobra.Command, args []string) error {
	npmCommand := "npm"
	if runtime.GOOS == "windows" {
		npmCommand = "npm.cmd"
	}

	initCmd := exec.Command(npmCommand, "init", "playwright@latest")
	initCmd.Stdin = os.Stdin
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr

	if err := initCmd.Run(); err != nil {
		return fmt.Errorf("run npm init playwright@latest error: %w", err)
	}

	return nil
}
