package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const defaultNvmVersion = "v0.40.3"

type NvmInstallCmd struct {
	Version string
	Force   bool
}

func (c *NvmInstallCmd) Run(cmd *cobra.Command, args []string) error {
	return installNvm(c.Version, c.Force)
}

func installNvm(version string, force bool) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("nvm install is not supported on windows")
	}

	if version == "" {
		version = defaultNvmVersion
	}

	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	nvmDir := filepath.Join(userHomePath, ".nvm")
	nvmShPath := filepath.Join(nvmDir, "nvm.sh")
	if !force {
		if _, err := os.Stat(nvmShPath); err == nil {
			return nil
		}
	}

	if force {
		os.RemoveAll(nvmDir)
	}

	scriptPath := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("nvm-install-%s.sh", strings.TrimPrefix(version, "v")),
	)
	os.Remove(scriptPath)

	downloadURL := fmt.Sprintf(
		"https://raw.githubusercontent.com/nvm-sh/nvm/%s/install.sh",
		version,
	)
	if err := downloadToFile(downloadURL, scriptPath); err != nil {
		return fmt.Errorf("download nvm install script error: %w", err)
	}
	defer os.Remove(scriptPath)

	installCmd := exec.Command("bash", scriptPath)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	installCmd.Env = append(
		os.Environ(),
		fmt.Sprintf("NVM_DIR=%s", nvmDir),
	)
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("run nvm install script error: %w", err)
	}

	if _, err := os.Stat(nvmShPath); err != nil {
		return fmt.Errorf("nvm install seems incomplete: %w", err)
	}

	return nil
}
