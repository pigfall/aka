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

type PlaywrightRunTestCmd struct {
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

func (c *PlaywrightRunTestCmd) Run(cmd *cobra.Command, args []string) error {
	npxCommand := "npx"
	if runtime.GOOS == "windows" {
		npxCommand = "npx.cmd"
	}

	testCmd := exec.Command(npxCommand, "playwright", "test")
	testCmd.Stdin = os.Stdin
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr

	if err := testCmd.Run(); err != nil {
		return fmt.Errorf("run npx playwright test error: %w", err)
	}

	return nil
}
