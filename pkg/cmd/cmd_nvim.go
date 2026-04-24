package cmd

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"

    "github.com/spf13/cobra"
)

// Uses pure-Go archive unpacking where possible; git and nvim exec calls remain

type NvimInstallCmd struct {
	InstallPlugin       bool
	InstallNodeJSForCoC bool
	NodeJSVersionForCoC string
}

func (n *NvimInstallCmd) Run(cmd *cobra.Command, args []string) error {
	if !n.InstallNodeJSForCoC {
		n.NodeJSVersionForCoC = ""
	}
	return nvimInstall(n.InstallPlugin, n.NodeJSVersionForCoC)
}

func nvimInstall(installPlugin bool, nodejsVersionForCoC string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	urls := map[string]string{
		"linux-amd64":   "https://github.com/neovim/neovim/releases/download/v0.12.2/nvim-linux-x86_64.tar.gz",
		"linux-arm64":   "https://github.com/neovim/neovim/releases/download/v0.12.2/nvim-linux-arm64.tar.gz",
		"darwin-arm64":  "https://github.com/neovim/neovim/releases/download/v0.12.2/nvim-macos-arm64.tar.gz",
		"darwin-amd64":  "https://github.com/neovim/neovim/releases/download/v0.12.2/nvim-macos-x86_64.tar.gz",
		"windows-amd64": "https://github.com/neovim/neovim/releases/download/v0.12.2/nvim-win64.zip",
		"windows-arm64": "https://github.com/neovim/neovim/releases/download/v0.12.2/nvim-win64.zip", // Not officially provided, fallback to amd64 build
	}
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	downloadURL := urls[platform]
	if downloadURL == "" {
		return fmt.Errorf("unsupported platform %s", platform)
	}

	installPath := filepath.Join(homeDir, "tools", "nvim")
	os.RemoveAll(installPath)
	os.MkdirAll(installPath, 0755)

    os.Remove("nvim.tar.gz")
    if err := downloadToFile(downloadURL, "nvim.tar.gz"); err != nil {
        return err
    }
    if err := UnpackArchive("nvim.tar.gz", installPath, 1); err != nil {
        return err
    }

	shrc := []string{
		".bashrc",
		".zshrc",
	}
	export := "\nexport PATH=$PATH:" + installPath + `/bin`
	for _, v := range shrc {
		v = filepath.Join(homeDir, v)
		if _, err := os.Stat(v); err == nil {
			f, err := os.OpenFile(v, os.O_RDWR|os.O_APPEND, os.ModePerm)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := f.WriteString(export); err != nil {
				return err
			}
		}
	}

	if installPlugin {
		if nodejsVersionForCoC != "" {
			_, err := installNodejs(nodejsVersionForCoC, "nodejs-coc")
			if err != nil {
				return fmt.Errorf("install nodejs for coc: %w", err)
			}
		}
		// install ripgrep.
		if err := installRipgrep(false); err != nil {
			return fmt.Errorf("install ripgrep: %w", err)
		}
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
