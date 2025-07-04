package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
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
	downloadPath := filepath.Join(downloadDir(), "kubectl")
	saveTo, err := os.Create(downloadPath)
	c.FailOnError(err)
	defer saveTo.Close()

	c.FailOnError(
		download(
			fmt.Sprintf("https://dl.k8s.io/release/%s/bin/linux/amd64/kubectl", version),
			saveTo,
		),
	)

	finalPath := filepath.Join(toolDir(), "kubectl")
	c.FailOnError(os.Rename(downloadPath, finalPath))
	c.FailOnError(os.Chmod(finalPath, 0744))

	return nil
}
