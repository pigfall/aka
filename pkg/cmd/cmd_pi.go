package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// defaultPiVersion is the Pi Coding Agent release installed when no --version
// is given. Pi publishes standalone platform archives on GitHub releases, so
// no Node.js runtime is required.
const defaultPiVersion = "1.6.0"

type PiInstallCmd struct {
	Version string
	Force   bool
}

func (c *PiInstallCmd) Run(cmd *cobra.Command, args []string) error {
	return installPi(c.Version, c.Force)
}

func installPi(version string, force bool) error {
	version = normalizePiVersion(version)
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	archiveName, stripComponents, err := piArchive(platform)
	if err != nil {
		return err
	}

	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dstPath := filepath.Join(userHomePath, "tools", "pi")
	binaryName := "pi"
	if runtime.GOOS == "windows" {
		binaryName = "pi.exe"
	}
	binaryPath := filepath.Join(dstPath, binaryName)

	if force {
		if err := os.RemoveAll(dstPath); err != nil {
			return fmt.Errorf("remove %s error: %w", dstPath, err)
		}
	} else if _, err := os.Stat(binaryPath); err == nil {
		// Already installed; keep the existing version instead of
		// re-downloading the same archive.
		return nil
	}

	downloadURL := fmt.Sprintf(
		"https://github.com/earendil-works/pi/releases/download/v%s/%s",
		version,
		archiveName,
	)
	downloadPath := filepath.Join(downloadDir(), archiveName)
	os.Remove(downloadPath)
	defer os.Remove(downloadPath)

	if err := downloadToFile(downloadURL, downloadPath); err != nil {
		return fmt.Errorf("download from %s error: %w", downloadURL, err)
	}

	if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
		return fmt.Errorf("mkdir %s error: %w", dstPath, err)
	}

	// The unix archives wrap everything in a top-level "pi/" directory while
	// the windows zip stores files at the archive root.
	if err := UnpackArchive(downloadPath, dstPath, stripComponents); err != nil {
		return fmt.Errorf("uncompress %s error: %w", downloadPath, err)
	}

	if _, err := os.Stat(binaryPath); err != nil {
		return fmt.Errorf("pi binary not found at %s", binaryPath)
	}
	if runtime.GOOS != "windows" {
		if err := os.Chmod(binaryPath, 0755); err != nil {
			return fmt.Errorf("chmod %s error: %w", binaryPath, err)
		}
	}

	if err := appendPathToUserShellRc(dstPath); err != nil {
		return err
	}

	return nil
}

// appendPathToUserShellRc prepends binDir to PATH in the shell rc files that
// already exist under the user's home directory. The export line is written
// only once, so repeated installs do not pile up duplicate entries.
func appendPathToUserShellRc(binDir string) error {
	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	exportLine := fmt.Sprintf("export PATH=%s:$PATH", binDir)
	for _, name := range []string{".bashrc", ".zshrc"} {
		rcPath := filepath.Join(userHomePath, name)
		content, err := os.ReadFile(rcPath)
		if err != nil {
			// rc file does not exist (or is unreadable): skip it.
			continue
		}
		if strings.Contains(string(content), exportLine) {
			continue
		}

		f, err := os.OpenFile(rcPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		if _, err := f.WriteString("\n" + exportLine + "\n"); err != nil {
			f.Close()
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}

	return nil
}

func normalizePiVersion(version string) string {
	version = strings.TrimSpace(version)
	if version == "" {
		return defaultPiVersion
	}
	return strings.TrimPrefix(version, "v")
}

// piArchive maps a GOOS-GOARCH pair to the Pi release archive name and the
// number of leading path elements to strip while unpacking.
func piArchive(platform string) (archiveName string, stripComponents int, err error) {
	switch platform {
	case "darwin-arm64":
		return "pi-darwin-arm64.tar.gz", 1, nil
	case "darwin-amd64":
		return "pi-darwin-x64.tar.gz", 1, nil
	case "linux-arm64":
		return "pi-linux-arm64.tar.gz", 1, nil
	case "linux-amd64":
		return "pi-linux-x64.tar.gz", 1, nil
	case "windows-arm64":
		return "pi-windows-arm64.zip", 0, nil
	case "windows-amd64":
		return "pi-windows-x64.zip", 0, nil
	default:
		return "", 0, fmt.Errorf("unsupported platform: %s", platform)
	}
}
