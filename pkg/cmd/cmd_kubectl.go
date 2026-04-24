package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type KubectlInstallCmd struct {
	Cmd
	Version string
}

func (c *KubectlInstallCmd) Run(cmd *cobra.Command, args []string) error {
	version := c.Version
	if version == "" {
		s := strings.Builder{}
		if err := download(
			"https://dl.k8s.io/release/stable.txt",
			&s,
		); err != nil {
			return fmt.Errorf("query latest kubeclt version error: %w", err)
		}
		version = s.String()
	}

	platform := runtime.GOOS
	arch := runtime.GOARCH

	var urlTpl string
	switch platform {
	case "linux":
		switch arch {
		case "amd64":
			urlTpl = "https://dl.k8s.io/release/%s/bin/linux/amd64/kubectl"
		case "arm64":
			urlTpl = "https://dl.k8s.io/release/%s/bin/linux/arm64/kubectl"
		default:
			panic(fmt.Sprintf("Unsupported arch for linux: %s", arch))
		}
	case "darwin":
		switch arch {
		case "amd64":
			urlTpl = "https://dl.k8s.io/release/%s/bin/darwin/amd64/kubectl"
		case "arm64":
			urlTpl = "https://dl.k8s.io/release/%s/bin/darwin/arm64/kubectl"
		default:
			panic(fmt.Sprintf("Unsupported arch for darwin: %s", arch))
		}
	case "windows":
		switch arch {
		case "amd64":
			urlTpl = "https://dl.k8s.io/release/%s/bin/windows/amd64/kubectl.exe"
		case "arm64":
			urlTpl = "https://dl.k8s.io/release/%s/bin/windows/arm64/kubectl.exe"
		default:
			panic(fmt.Sprintf("Unsupported arch for windows: %s", arch))
		}
	default:
		panic(fmt.Sprintf("Unsupported OS: %s", platform))
	}

	downloadPath := filepath.Join(downloadDir(), "kubectl")
	saveTo, err := os.Create(downloadPath)
	c.FailOnError(err)
	defer saveTo.Close()

	c.FailOnError(
		download(
			fmt.Sprintf(urlTpl, version),
			saveTo,
		),
	)

	finalPath := filepath.Join(toolDir(), "kubectl")
	c.FailOnError(os.Rename(downloadPath, finalPath))
	c.FailOnError(os.Chmod(finalPath, 0744))

	return nil
}
