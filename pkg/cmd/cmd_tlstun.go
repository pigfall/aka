package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

type TLSTunInstallCmd struct {
}

func (c *TLSTunInstallCmd) Run(cmd *cobra.Command, args []string) error {
	tmpDir := os.TempDir()
	sourcePath := filepath.Join(tmpDir, "tlstun")
	os.RemoveAll(sourcePath)

	gitClone := exec.Command("git", "clone", "--depth=1", "https://github.com/pigfall/tlstun", sourcePath)
	gitClone.Stderr = os.Stderr
	gitClone.Stdout = os.Stdout
	if err := gitClone.Run(); err != nil {
		return fmt.Errorf("git clone error: %w", err)
	}

	install := exec.Command("go", "install", ".")
	install.Dir = sourcePath
	install.Stderr = os.Stderr
	install.Stdout = os.Stdout
	if err := install.Run(); err != nil {
		return fmt.Errorf("go install error: %w", err)
	}

	return nil
}
