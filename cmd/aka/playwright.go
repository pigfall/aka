package main

import (
	pkgcmd "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func playwrightCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "playwright",
	}

	initPlaywright := pkgcmd.PlaywrightInitCmd{}
	initPlaywrightCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize Playwright in current folder",
		Args:  cobra.NoArgs,
		RunE:  initPlaywright.Run,
	}

	runPlaywrightTest := pkgcmd.PlaywrightRunTestCmd{}
	runPlaywrightTestCmd := &cobra.Command{
		Use:   "run-test",
		Short: "Run Playwright tests in current folder",
		Args:  cobra.NoArgs,
		RunE:  runPlaywrightTest.Run,
	}

	c.AddCommand(initPlaywrightCmd, runPlaywrightTestCmd)

	return c
}
