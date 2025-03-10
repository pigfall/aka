package cmd

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	toml "github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

//go:embed assets/*
var assets embed.FS

const craftingSandboxSnapshotIncludes = `
.vscode-remote/extensions
.snapshot.personal
.gitconfig
.config/nvim
.sandbox.personal/bashrc
tools
.vim
`

const vscodeShortcutJSONConfig = `
[
{
    "key":"ctrl+o",
    "command":"editor.action.triggerSuggest",
    "when": "editorTextFocus && vim.active && vim.mode=='Insert'"
},
{
    "key": "ctrl+e",
    "command": "cursorEnd",
    "when": "vim.active && vim.mode=='Insert'"
}
]
`

const kittyConfig = `
font_size 14.0

map alt+1 goto_tab 1
map alt+2 goto_tab 2
map alt+3 goto_tab 3
map alt+4 goto_tab 4
map alt+5 goto_tab 5
map alt+6 goto_tab 6

map ctrl+c send_text all \x03
`

// Personalize git.
type PersonalizeGitCmd struct {
}

// Personalize crafting sandbox.
type PersonalizeCraftingSandboxCmd struct {
	SnapshotName string
}

// Personalize vscode.
type PersonalizeVscodeCmd struct {
}

// Personalize kitty terminal
type PersonalizeKittyCmd struct{}

// Personalize tlstun
type PersonalizeTLSTun struct {
	Password string
}

func (c *PersonalizeTLSTun) Run(cmd *cobra.Command, args []string) error {
	return personalizeTLSTun(c.Password)
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

func (c *PersonalizeVscodeCmd) Run(cobraCmd *cobra.Command, args []string) error {
	var configLocation string
	userConfigPath, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("get user config directory error: %w", err)
	}
	switch ostype := runtime.GOOS; ostype {
	case "darwin":
		configLocation = filepath.Join(userConfigPath, "Code", "User", "keybindings.json")
	case "windows":
		configLocation = filepath.Join(userConfigPath, "Code", "User", "keybindings.json")
	case "linux":
		configLocation = filepath.Join(userConfigPath, "Code", "User", "keybindings.json")
	default:
		return fmt.Errorf("unsupported platform: %s", ostype)
	}

	if err := os.WriteFile(configLocation, []byte(vscodeShortcutJSONConfig), os.ModePerm); err != nil {
		return fmt.Errorf("write to %s error: %w", configLocation, err)
	}

	return nil
}

func (c *PersonalizeKittyCmd) Run(cobraCmd *cobra.Command, args []string) error {
	var configPath string
	switch ostype := runtime.GOOS; ostype {
	case "darwin", "linux":
		userHomePath, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get user home directory error: %w", err)
		}
		configPath = filepath.Join(userHomePath, ".config/kitty/kitty.conf")
	case "windows":
		userConfigPath, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("get user config directory error: %w", err)
		}
		configPath = filepath.Join(userConfigPath, "kitty", "kitty.conf")
	default:
		return fmt.Errorf("unsupported platform: %s", ostype)
	}
	os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
	if err := os.WriteFile(configPath, []byte(kittyConfig), os.ModePerm); err != nil {
		return fmt.Errorf("write to %s error: %w", configPath, err)
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

func personalizeTLSTun(password string) error {
	userHomePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	ca, err := fs.ReadFile(assets, "assets/acnutslt")
	if err != nil {
		return fmt.Errorf("read ca error: %w", err)
	}
	clientCert, err := fs.ReadFile(assets, "assets/clicnutslt")
	if err != nil {
		return fmt.Errorf("read client cert error: %w", err)
	}
	clientKey, err := fs.ReadFile(assets, "assets/cliknutslt")
	if err != nil {
		return fmt.Errorf("read client key error: %w", err)
	}

	ca, err = decrypt(ca, password)
	if err != nil {
		return fmt.Errorf("decrypt ca error: %w", err)
	}

	clientCert, err = decrypt(clientCert, password)
	if err != nil {
		return fmt.Errorf("decrypt client cert error: %w", err)
	}

	clientKey, err = decrypt(clientKey, password)
	if err != nil {
		return fmt.Errorf("decrypt client key error: %w", err)
	}

	caPath := filepath.Join(userHomePath, "tlstun-ca.pem")
	if err := os.WriteFile(caPath, ca, os.ModePerm); err != nil {
		return fmt.Errorf("write to %s error: %w", caPath, err)
	}
	clientCertPath := filepath.Join(userHomePath, "tlstun-clientcert.pem")
	if err := os.WriteFile(clientCertPath, clientCert, os.ModePerm); err != nil {
		return fmt.Errorf("write to %s error: %w", clientCertPath, err)
	}
	clientKeyPath := filepath.Join(userHomePath, "tlstun-clientkey.pem")
	if err := os.WriteFile(clientKeyPath, clientKey, os.ModePerm); err != nil {
		return fmt.Errorf("write to %s error: %w", clientKeyPath, err)
	}

	serverCert, err := fs.ReadFile(assets, "assets/tsltun-servercert.pem")
	if err != nil {
		return err
	}
	serverKey, err := fs.ReadFile(assets, "assets/tlstun-serverkey.pem")
	if err != nil {
		return err
	}

	serverCertPath := filepath.Join(userHomePath, "tsltun-servercert.pem")
	if err := os.WriteFile(serverCertPath, serverCert, os.ModePerm); err != nil {
		return fmt.Errorf("write to %s error: %w", serverCertPath, err)
	}
	serverKeyPath := filepath.Join(userHomePath, "tlstun-serverkey.pem")
	if err := os.WriteFile(serverKeyPath, serverKey, os.ModePerm); err != nil {
		return fmt.Errorf("write to %s error: %w", serverKeyPath, err)
	}

	return nil
}
