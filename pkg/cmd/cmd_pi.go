package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const piReleasesBaseURL = "https://github.com/earendil-works/pi/releases"

// piLatestReleaseURL is the "latest release" endpoint. GitHub answers it with a
// redirect to /releases/tag/vX.Y.Z, so aka can resolve the newest version
// without consuming GitHub API rate limit. It is a variable to keep the
// resolver testable.
var piLatestReleaseURL = piReleasesBaseURL + "/latest"

type PiInstallCmd struct {
	Version string
	Force   bool
}

func (c *PiInstallCmd) Run(cmd *cobra.Command, args []string) error {
	return installPi(c.Version, c.Force)
}

func installPi(version string, force bool) error {
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	archiveName, stripComponents, err := piArchive(platform)
	if err != nil {
		return err
	}

	version = normalizePiVersion(version)
	if version == "" {
		latest, err := latestPiVersion()
		if err != nil {
			return err
		}
		version = latest
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

	if !force {
		// Skip when the requested release is already on disk; a plain
		// `aka pi install` therefore upgrades to the latest release only
		// when the installed one is older.
		if installedPiVersion(dstPath) == version {
			if _, err := os.Stat(binaryPath); err == nil {
				return nil
			}
		}
	}
	if force {
		if err := os.RemoveAll(dstPath); err != nil {
			return fmt.Errorf("remove %s error: %w", dstPath, err)
		}
	}

	fmt.Fprintf(os.Stderr, "Installing pi v%s\n", version)

	downloadURL := fmt.Sprintf(
		"%s/download/v%s/%s",
		piReleasesBaseURL,
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

// latestPiVersion resolves the newest Pi release tag, e.g. "1.6.0".
func latestPiVersion() (string, error) {
	client, err := httpClientFromEnv()
	if err != nil {
		return "", fmt.Errorf("creating http client: %w", err)
	}

	// Copy the client so the shared instance keeps its redirect behaviour.
	noRedirect := *client
	noRedirect.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest(http.MethodGet, piLatestReleaseURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request for %s error: %w", piLatestReleaseURL, err)
	}
	resp, err := noRedirect.Do(req)
	if err != nil {
		return "", fmt.Errorf("resolve latest pi release from %s error: %w", piLatestReleaseURL, err)
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("resolve latest pi release from %s: unexpected status %s", piLatestReleaseURL, resp.Status)
	}

	u, err := url.Parse(location)
	if err != nil {
		return "", fmt.Errorf("parse release location %s error: %w", location, err)
	}
	version := normalizePiVersion(path.Base(u.Path))
	if version == "" {
		return "", fmt.Errorf("resolve latest pi release from %s: no version in %s", piLatestReleaseURL, location)
	}

	return version, nil
}

// installedPiVersion reads the version recorded in the installed archive's
// package.json. It returns "" when pi is not installed or the metadata is
// unreadable, which makes the caller reinstall.
func installedPiVersion(piDir string) string {
	content, err := os.ReadFile(filepath.Join(piDir, "package.json"))
	if err != nil {
		return ""
	}

	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return ""
	}

	return strings.TrimSpace(pkg.Version)
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

// normalizePiVersion trims spaces and a leading "v". An empty result means
// "resolve the latest release".
func normalizePiVersion(version string) string {
	return strings.TrimPrefix(strings.TrimSpace(version), "v")
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
