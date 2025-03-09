package main

import (
	pkgcmd "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func personalizeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "personalize",
		Aliases: []string{"pl"},
	}

	git := pkgcmd.PersonalizeGitCmd{}
	gitCmd := &cobra.Command{
		Use:  "git",
		RunE: git.Run,
	}

	craftingSandbox := pkgcmd.PersonalizeCraftingSandboxCmd{}
	craftingSandboxCmd := &cobra.Command{
		Use:  "crafting-sandbox",
		RunE: craftingSandbox.Run,
	}
	craftingSandboxCmd.Flags().StringVar(
		&craftingSandbox.SnapshotName,
		"snapshot-name",
		"tzz-personal",
		"snapshot name",
	)

	cmd.AddCommand(
		gitCmd,
		craftingSandboxCmd,
	)

	return cmd
}
