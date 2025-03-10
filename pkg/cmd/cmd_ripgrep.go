package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

type RipgrepInstallCmd struct {
	Force bool
}

func (c *RipgrepInstallCmd) Run(cmd *cobra.Command, args []string) error {
	return installRipgrep(c.Force)
}

func installRipgrep(force bool) error {
	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dstPath := filepath.Join(userHomePath, "tools", "ripgrep")
	if force {
		os.RemoveAll(dstPath)
	}

	ripgrepBinPath := filepath.Join(dstPath, "rg")
	if _, err := os.Stat(ripgrepBinPath); err == nil {
		return nil
	}

	var urls = map[string]string{
		"darwin-arm64": "https://github.com/BurntSushi/ripgrep/releases/download/14.1.1/ripgrep-14.1.1-aarch64-apple-darwin.tar.gz",
		"linux-amd64":  "https://github.com/BurntSushi/ripgrep/releases/download/14.1.1/ripgrep-14.1.1-x86_64-unknown-linux-musl.tar.gz",
	}

	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	url := urls[platform]
	if url == "" {
		return fmt.Errorf("unsupported platform: %s", platform)
	}

	downloadPath := filepath.Join(userHomePath, "tools", "ripgrep.tar.gz")
	os.MkdirAll(filepath.Dir(downloadPath), os.ModePerm)

	downloadCmd := exec.Command(
		"curl",
		"-o",
		downloadPath,
		"-L",
		url,
	)
	downloadCmd.Stdout = os.Stdout
	downloadCmd.Stderr = os.Stderr
	if err := downloadCmd.Run(); err != nil {
		return fmt.Errorf("download from %s error: %w", url, err)
	}

	os.MkdirAll(dstPath, os.ModePerm)
	uncompressCmd := exec.Command("tar", "-xf", downloadPath, "--strip-components=1", "-C", dstPath)
	uncompressCmd.Stdout = os.Stdout
	uncompressCmd.Stderr = os.Stderr
	if err := uncompressCmd.Run(); err != nil {
		return fmt.Errorf("uncompress %s error: %w", downloadPath, err)
	}

	shrc := []string{
		".bashrc",
		".zshrc",
	}
	for _, v := range shrc {
		v = filepath.Join(userHomePath, v)
		if _, err := os.Stat(v); err == nil {
			f, err := os.OpenFile(v, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := f.WriteString(fmt.Sprintf("\nexport PATH=$PATH:%s", dstPath)); err != nil {
				return err
			}
		}
	}

	return nil
}
