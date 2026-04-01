package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

const installGHBashScript = `
set -e
set -o pipefail

sudo mkdir -p -m 755 /etc/apt/keyrings
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null
sudo chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt-get update -y
sudo apt-get install gh -y
`

type InstallGHCmd struct{}

func (c *InstallGHCmd) Run(cobraCmd *cobra.Command, args []string) error {
	ctx := cobraCmd.Context()
	cmd := exec.CommandContext(
		ctx,
		"bash",
		"-s",
	)
	cmd.Stdin = strings.NewReader(installGHBashScript)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	return nil
}
