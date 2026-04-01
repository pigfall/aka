package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

const defaultGHVersion = "2.88.1"

type InstallGHCmd struct {
	Version string
}

func (c *InstallGHCmd) Run(cmd *cobra.Command, args []string) error {
	return installGH(c.Version)
}

func installGH(version string) error {
	if version == "" {
		version = defaultGHVersion
	}

	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dstPath := filepath.Join(userHomePath, "tools", "gh")
	os.RemoveAll(dstPath)
	os.MkdirAll(dstPath, os.ModePerm)

	var downloadURLBuilders = map[string]func(version string) string{
		"linux-amd64": func(version string) string {
			return fmt.Sprintf("https://github.com/cli/cli/releases/download/v%s/gh_%s_linux_amd64.tar.gz", version, version)
		},
		"linux-arm64": func(version string) string {
			return fmt.Sprintf("https://github.com/cli/cli/releases/download/v%s/gh_%s_linux_arm64.tar.gz", version, version)
		},
		"darwin-amd64": func(version string) string {
			return fmt.Sprintf("https://github.com/cli/cli/releases/download/v%s/gh_%s_macOS_amd64.zip", version, version)
		},
		"darwin-arm64": func(version string) string {
			return fmt.Sprintf("https://github.com/cli/cli/releases/download/v%s/gh_%s_macOS_arm64.zip", version, version)
		},
		"windows-amd64": func(version string) string {
			return fmt.Sprintf("https://github.com/cli/cli/releases/download/v%s/gh_%s_windows_amd64.zip", version, version)
		},
		"windows-arm64": func(version string) string {
			return fmt.Sprintf("https://github.com/cli/cli/releases/download/v%s/gh_%s_windows_arm64.zip", version, version)
		},
	}

	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	urlBuilder := downloadURLBuilders[platform]
	if urlBuilder == nil {
		return fmt.Errorf("unsupported platform: %s", platform)
	}
	downloadURL := urlBuilder(version)

	ext := ".tar.gz"
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		ext = ".zip"
	}

	downloadPath := filepath.Join(os.TempDir(), "gh"+ext)
	os.Remove(downloadPath)

	downloadCmd := exec.Command("curl", "-L", "-o", downloadPath, downloadURL)
	downloadCmd.Stdout = os.Stdout
	downloadCmd.Stderr = os.Stderr
	if err := downloadCmd.Run(); err != nil {
		return fmt.Errorf("download from %s error: %w", downloadURL, err)
	}

	if ext == ".tar.gz" {
		uncompressCmd := exec.Command("tar", "-xf", downloadPath, "--strip-components=1", "-C", dstPath)
		uncompressCmd.Stdout = os.Stdout
		uncompressCmd.Stderr = os.Stderr
		if err := uncompressCmd.Run(); err != nil {
			return fmt.Errorf("uncompress %s error: %w", downloadPath, err)
		}
	} else {
		uncompressCmd := exec.Command("unzip", "-o", downloadPath, "-d", dstPath)
		uncompressCmd.Stdout = os.Stdout
		uncompressCmd.Stderr = os.Stderr
		if err := uncompressCmd.Run(); err != nil {
			return fmt.Errorf("uncompress %s error: %w", downloadPath, err)
		}
		// zip extracts into a subdirectory; move contents up
		entries, err := os.ReadDir(dstPath)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				subDir := filepath.Join(dstPath, entry.Name())
				subEntries, err := os.ReadDir(subDir)
				if err != nil {
					return err
				}
				for _, se := range subEntries {
					if err := os.Rename(
						filepath.Join(subDir, se.Name()),
						filepath.Join(dstPath, se.Name()),
					); err != nil {
						return err
					}
				}
				os.RemoveAll(subDir)
				break
			}
		}
	}

	binDir := filepath.Join(dstPath, "bin")
	shrc := []string{
		".bashrc",
		".zshrc",
	}
	for _, v := range shrc {
		v = filepath.Join(userHomePath, v)
		if _, err := os.Stat(v); err == nil {
			f, err := os.OpenFile(v, os.O_RDWR|os.O_APPEND, os.ModePerm)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := f.WriteString(fmt.Sprintf("\nexport PATH=$PATH:%s", binDir)); err != nil {
				return err
			}
		}
	}

	return nil
}
