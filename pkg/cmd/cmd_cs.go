package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type CsLoginProviderAsanaCreateCmd struct {
	Cmd
	OrgName      string
	AuthURL      string
	TokenURL     string
	ClientID     string
	ClientSecret string
}

type loginProvdierAsana struct {
	Provider loginProviderAsanaProvider `json:"provider"`
	Policy   loginProviderAsanaPolicy   `json:"policy"`
}

type loginProviderAsanaProvider struct {
	OAuth2 loginProviderAsanaProviderOAuth2 `json:"oauth2"`
}

type loginProviderAsanaProviderOAuth2 struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AuthURL      string `json:"auth_url"`
	TokenURL     string `json:"token_url"`
}

type loginProviderAsanaPolicy struct {
	AutoRefresh loginProviderAsanaPolicyAutoRefresh `json:"auto_refresh_token,omitempty"`
}

type loginProviderAsanaPolicyAutoRefresh struct {
	Enabled  bool `json:"enabled,omitempty"`
	MaxCount int  `json:"max_count,omitempty"`
}

type csSystemInfo struct {
	SysBaseURL string `json:"sys_base_url"`
}

func (c *CsLoginProviderAsanaCreateCmd) Run(cobraCmd *cobra.Command, args []string) error {
	name := args[0]

	infoCmd := exec.CommandContext(
		cobraCmd.Context(),
		"cs",
		"info",
		"--system",
		"--output-format=json",
		"--org="+c.OrgName,
	)
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	infoCmd.Stdout = &stdout
	infoCmd.Stderr = &stderr
	infoCmd.Stdin = nil

	if err := infoCmd.Run(); err != nil {
		return fmt.Errorf("run cs info: %w, stderr: %s", err, stderr.String())
	}
	info := csSystemInfo{}
	jsonDecoder := json.NewDecoder(strings.NewReader(stdout.String()))
	if err := jsonDecoder.Decode(&info); err != nil {
		return fmt.Errorf("decode cs info: %w", err)
	}
	if info.SysBaseURL == "" {
		return fmt.Errorf("cs info: missing sys_base_url")
	}

	callbackURL := strings.TrimSuffix(info.SysBaseURL, "/") + "/integration/login/callback"
	fmt.Printf("Add the following URL to your Asana app's OAuth callback URL: \x1b[32m%s\x1b[0m\n", callbackURL)

	fmt.Print("Press Enter to continue after adding the callback URL, or type 'n' to abort: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("read confirmation: %w", err)
	}
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "n" || input == "no" {
		return fmt.Errorf("aborted by user")
	}

	loginProvider := loginProvdierAsana{
		Provider: loginProviderAsanaProvider{
			OAuth2: loginProviderAsanaProviderOAuth2{
				ClientID:     c.ClientID,
				ClientSecret: c.ClientSecret,
				AuthURL:      c.AuthURL,
				TokenURL:     c.TokenURL,
			},
		},
		Policy: loginProviderAsanaPolicy{AutoRefresh: loginProviderAsanaPolicyAutoRefresh{Enabled: true, MaxCount: 3}},
	}

	// Execute cs org login-provider create <NAME> --from=-
	createCmd := exec.CommandContext(
		cobraCmd.Context(),
		"cs",
		"org",
		"login-provider",
		"create",
		name,
		"--from=-",
		"--org="+c.OrgName,
	)
	createCmd.Stderr = os.Stderr
	createCmd.Stdout = os.Stdout
	jsonData, err := json.Marshal(loginProvider)
	if err != nil {
		return fmt.Errorf("marshal login provider: %w", err)
	}
	createCmd.Stdin = strings.NewReader(string(jsonData))

	return createCmd.Run()
}
