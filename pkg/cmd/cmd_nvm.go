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

type NvmInstallCmd struct {
}

func (c *NvmInstallCmd) Run(cmd *cobra.Command, args []string) error {
	return installNvm("v0.40.3", false)
}

type NvmNodejsInstallCmd struct {
}

func (c *NvmNodejsInstallCmd) Run(cmd *cobra.Command, args []string) error {
	return installNodejsByNvm(args[0])
}

type NvmListCmd struct {
}

func (c *NvmListCmd) Run(cmd *cobra.Command, args []string) error {
	return listNodejsByNvm()
}

func installNvm(version string, force bool) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("nvm install is not supported on windows")
	}

	if version == "" {
		version = "v0.40.3"
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

func installNodejsByNvm(version string) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("nvm install nodejs is not supported on windows")
	}

	installCmd := exec.Command("bash", "-c", `
set -e
NVM_DIR="${NVM_DIR:-$HOME/.nvm}"
if [ ! -s "$NVM_DIR/nvm.sh" ]; then
	echo "nvm is not installed. please install nvm first." >&2
	exit 1
fi
. "$NVM_DIR/nvm.sh"
nvm install "$AKA_NODEJS_VERSION"
`)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	installCmd.Env = append(installCmd.Env, fmt.Sprintf("AKA_NODEJS_VERSION=%s", version))
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("run nvm install %s error: %w", version, err)
	}

	return nil
}

func listNodejsByNvm() error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("nvm list is not supported on windows")
	}

	listCmd := exec.Command("bash", "-c", `
set -e
NVM_DIR="${NVM_DIR:-$HOME/.nvm}"
if [ ! -s "$NVM_DIR/nvm.sh" ]; then
	echo "nvm is not installed. please install nvm first." >&2
	exit 1
fi
. "$NVM_DIR/nvm.sh"
nvm list
`)
	listCmd.Stdout = os.Stdout
	listCmd.Stderr = os.Stderr
	listCmd.Env = append(listCmd.Env, os.Environ()...)
	if err := listCmd.Run(); err != nil {
		return fmt.Errorf("run nvm list error: %w", err)
	}

	return nil
}
