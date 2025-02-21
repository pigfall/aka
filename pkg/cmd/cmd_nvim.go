package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

type NvimInstallCmd struct {
	InstallPlugin bool
}

func (n *NvimInstallCmd) Run(cmd *cobra.Command, args []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	installPath := filepath.Join(homeDir, "tools", "nvim")
	os.RemoveAll(installPath)
	os.MkdirAll(installPath, 0755)

	download := exec.Command("curl", "-L", "-o", "nvim.tar.gz", "https://github.com/neovim/neovim/releases/download/v0.10.4/nvim-linux-x86_64.tar.gz")
	download.Stdout = os.Stdout
	download.Stderr = os.Stderr
	if err := download.Run(); err != nil {
		return err
	}

	install := exec.Command("tar", "-xf", "nvim.tar.gz", "--strip-components=1", "-C", installPath)
	install.Stdout = os.Stdout
	install.Stderr = os.Stderr
	if err := install.Run(); err != nil {
		return err
	}

	export := `export PATH=$PATH:` + installPath + `/bin`
	bashrc, err := os.OpenFile(filepath.Join(homeDir, ".bashrc"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Open %s: %w", filepath.Join(homeDir, ".bashrc"), err)
	}
	if _, err := bashrc.Write([]byte(export)); err != nil {
		return fmt.Errorf("Update bashrc: %w", err)
	}
	bashrc.Close()

	if n.InstallPlugin {
		nvimPluginDir := filepath.Join(homeDir, ".config", "nvim")
		os.RemoveAll(nvimPluginDir)
		os.MkdirAll(filepath.Dir(nvimPluginDir), 0755)
		pluginDownload := exec.Command("git", "clone", "https://github.com/pigfall/nvimc2.git", nvimPluginDir)
		pluginDownload.Stdout = os.Stdout
		pluginDownload.Stderr = os.Stderr
		if err := pluginDownload.Run(); err != nil {
			return fmt.Errorf("Download nvim plugin: %w", err)
		}
		installPlugin := exec.Command(filepath.Join(installPath, "bin", "nvim"), "-u", filepath.Join(nvimPluginDir, "plugins.vim"), "--headless", "-c", "PlugInstall", "-c", "qa")
		installPlugin.Stdout = os.Stdout
		installPlugin.Stderr = os.Stderr
		if err := installPlugin.Run(); err != nil {
			return fmt.Errorf("Install nvim plugin: %w", err)
		}
	}

	return nil
}
