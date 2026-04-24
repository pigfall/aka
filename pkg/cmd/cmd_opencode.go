package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"

    "github.com/spf13/cobra"
)

const defaultOpenCodeVersion = "1.4.3"

type InstallOpenCodeCmd struct {
	Version string
}

func (c *InstallOpenCodeCmd) Run(cmd *cobra.Command, args []string) error {
	return installOpenCode(c.Version)
}

func installOpenCode(version string) error {
	if version == "" {
		version = defaultOpenCodeVersion
	}

	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dstPath := filepath.Join(userHomePath, "tools", "opencode")
	os.RemoveAll(dstPath)
	os.MkdirAll(dstPath, os.ModePerm)

	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	archiveName := opencodeArchiveName(platform)
	if archiveName == "" {
		return fmt.Errorf("unsupported platform: %s", platform)
	}

	downloadURL := fmt.Sprintf(
		"https://github.com/anomalyco/opencode/releases/download/v%s/%s",
		version,
		archiveName,
	)
	downloadPath := filepath.Join(os.TempDir(), archiveName)
	os.Remove(downloadPath)

	if err := downloadToFile(downloadURL, downloadPath); err != nil {
		return fmt.Errorf("download from %s error: %w", downloadURL, err)
	}

    // Use the pure-Go unpacker which supports .tar.gz and .zip
    if err := UnpackArchive(downloadPath, dstPath, 0); err != nil {
        return fmt.Errorf("uncompress %s error: %w", downloadPath, err)
    }

	binaryPath := filepath.Join(dstPath, "opencode")
	if runtime.GOOS == "windows" {
		binaryPath = filepath.Join(dstPath, "opencode.exe")
	}
	if _, err := os.Stat(binaryPath); err != nil {
		return fmt.Errorf("opencode binary not found at %s", binaryPath)
	}
	if runtime.GOOS != "windows" {
		if err := os.Chmod(binaryPath, 0755); err != nil {
			return fmt.Errorf("chmod %s error: %w", binaryPath, err)
		}
	}

	binDir := dstPath
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
			if _, err := f.WriteString(fmt.Sprintf("\nexport PATH=$PATH:%s", binDir)); err != nil {
				return err
			}
		}
	}

	return nil
}

func opencodeArchiveName(platform string) (archiveName string) {
	switch platform {
	case "darwin-arm64":
		archiveName = "opencode-darwin-arm64.zip"
	case "darwin-amd64":
		archiveName = "opencode-darwin-x64.zip"
	case "linux-arm64":
		archiveName = "opencode-linux-arm64.tar.gz"
	case "linux-amd64":
		archiveName = "opencode-linux-x64.tar.gz"
	case "windows-arm64":
		archiveName = "opencode-windows-arm64.zip"
	case "windows-amd64":
		archiveName = "opencode-windows-x64.zip"
	default:
		return ""
	}

	return archiveName
}
