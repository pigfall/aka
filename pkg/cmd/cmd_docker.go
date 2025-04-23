package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

const installDockerBashScirpt =`
set -e
set -o pipefail

# Add Docker's official GPG key:
sudo apt-get update -y
sudo apt-get install ca-certificates curl -y
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update -y

sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y

set +e
sudo groupadd docker
sudo usermod -aG docker $USER
`

type InstallDockerCmd struct{}

func (c *InstallDockerCmd)Run(cobraCmd *cobra.Command, args []string) error{
  ctx := cobraCmd.Context()
  cmd := exec.CommandContext(
      ctx,
      "bash",
      "-s",
  )
  cmd.Stdin = strings.NewReader(installDockerBashScirpt)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Env = append(os.Environ(),"DEBIAN_FRONTEND=noninteractive")
  if err :=cmd.Run();err != nil{
    os.Exit(1)
  }

  return nil
}
