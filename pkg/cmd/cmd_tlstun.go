package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

type TLSTunInstallCmd struct {
	Password string
}

type TLSTunClientCmd struct {
	TLSTunPath string
	Background bool
}

type TLSTunServerCmd struct{}

func (c *TLSTunClientCmd) Run(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("please provide server address")
	}
	serverAddr := args[0]
	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	cmd := exec.Command(
		c.TLSTunPath,
		"client",
		"-addr=127.0.0.1:1080",
		"-ca="+filepath.Join(userHomePath, "tlstun-ca.pem"),
		"-cert="+filepath.Join(userHomePath, "tlstun-clientcert.pem"),
		"-key="+filepath.Join(userHomePath, "tlstun-clientkey.pem"),
		serverAddr,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if c.Background {
		setRunInBackground(cmd)
		if err := cmd.Start(); err != nil {
			return err
		}
		cmd.Process.Release()
		return nil
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
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

	if c.Password != "" {
		if err := personalizeTLSTun(c.Password); err != nil {
			return err
		}
	}

	return nil
}
