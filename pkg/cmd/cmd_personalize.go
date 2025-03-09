package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

const craftingSandboxSnapshotIncludes = `
.vscode-remote/extensions
.snapshot.personal
.gitconfig
.config/nvim
.sandbox.personal/bashrc
tools
.vim
`

type PersonalizeGitCmd struct {
}

type PersonalizeCraftingSandboxCmd struct {
	SnapshotName string
}

func (c *PersonalizeGitCmd) Run(cmd *cobra.Command, args []string) error {
	return personalizeGit()
}

func (c *PersonalizeCraftingSandboxCmd) Run(cobraCmd *cobra.Command, args []string) error {
	if c.SnapshotName == "" {
		c.SnapshotName = "tzz-personal"
	}

	if err := personalizeGit(); err != nil {
		return fmt.Errorf("personalize gitconfig: %w", err)
	}
	if err := nvimInstall(true, ""); err != nil {
		return fmt.Errorf("install nvim: %w", err)
	}

	// download vim extension of vscode.
	vscodevimFilePath := "/tmp/vscodevim.vsix"
	os.Remove(vscodevimFilePath)
	cmd := exec.Command("curl", "-L", "-o", vscodevimFilePath, "https://openvsxorg.blob.core.windows.net/resources/vscodevim/vim/1.29.0/vscodevim.vim-1.29.0.vsix")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("download vim extension for vscode: %w", err)
	}
	// install vim extension.
	// /opt/sandboxd/vscode/bin/code-server-cs --install-extension ~/vscodevim.vim-1.29.0.vsix
	output, err := exec.Command("/opt/sandboxd/vscode/bin/code-server-cs", "--install-extension", vscodevimFilePath).CombinedOutput()
	if err != nil {
		return fmt.Errorf("install vim extension %w: %s", err, string(output))
	}

	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// update ~/.sandbox.personal/bashrc
	personalBashrcPath := filepath.Join(userHomePath, ".sandbox.personal", "bashrc")
	os.MkdirAll(filepath.Dir(personalBashrcPath), os.ModePerm)
	bashrc, err := os.OpenFile(personalBashrcPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open %s: %w", personalBashrcPath, err)
	}
	defer bashrc.Close()
	if _, err := bashrc.Write([]byte("\nexport PATH=$PATH:$HOME/tools/nvim/bin")); err != nil {
		return fmt.Errorf("write %s: %w", personalBashrcPath, err)
	}

	// create ~/snapshots.personal/includes.txt
	includesFilePath := filepath.Join(userHomePath, ".snapshot.personal", "includes.txt")
	if err := os.WriteFile(includesFilePath, []byte(craftingSandboxSnapshotIncludes), os.ModePerm); err != nil {
		return fmt.Errorf("update %s: %w", includesFilePath, err)
	}

	// create personal snapshot.
	createSnapshotCmd := exec.Command("cs", "snapshot", "create", c.SnapshotName, "--personal", "--force")
	createSnapshotCmd.Stdout = os.Stdout
	createSnapshotCmd.Stderr = os.Stderr
	if err := createSnapshotCmd.Run(); err != nil {
		return fmt.Errorf("create snapshot: %w", err)
	}

	return nil
}

func personalizeGit() error {
	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	gitConfigFilePath := filepath.Join(userHomePath, ".gitconfig")
	gitConfig := make(map[string]any)

	if _, err := os.Stat(gitConfigFilePath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		f, err := os.Create(gitConfigFilePath)
		if err != nil {
			return err
		}
		f.Close()
	} else {
		content, err := os.ReadFile(gitConfigFilePath)
		if err != nil {
			return err
		}
		if err := toml.Unmarshal(content, &gitConfig); err != nil {
			return err
		}
	}

	if gitConfig["core"] == nil {
		gitConfig["core"] = make(map[string]any)
	}
	gitConfig["core"].(map[string]any)["editor"] = "nvim"

	if gitConfig["alias"] == nil {
		gitConfig["alias"] = make(map[string]any)
	}
	gitConfig["alias"].(map[string]any)["lg"] = `log --graph --all --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit --date=relative`
	gitConfig["alias"].(map[string]any)["ss"] = `status`
	gitConfig["alias"].(map[string]any)["cm"] = `commit`
	gitConfig["alias"].(map[string]any)["ck"] = `checkout`

	if gitConfig["user"] == nil {
		gitConfig["user"] = make(map[string]any)
	}
	gitConfig["user"].(map[string]any)["name"] = "tzz"
	gitConfig["user"].(map[string]any)["email"] = "tangbe9@gmail.com"

	b, err := toml.Marshal(gitConfig)
	if err != nil {
		return err
	}

	return os.WriteFile(gitConfigFilePath, b, os.ModePerm)

}
