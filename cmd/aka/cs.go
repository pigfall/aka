package main

import (
	pkgcmd "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

const (
	asanaAuthURL  = "https://app.asana.com/-/oauth_authorize"
	asanaTokenURL = "https://app.asana.com/-/oauth_token"
)

func csCmd() *cobra.Command {
	asanaCreate := pkgcmd.CsLoginProviderAsanaCreateCmd{}
	createCmd := &cobra.Command{
		Use:  "create NAME [key=val...]",
		Args: cobra.MinimumNArgs(1),
		RunE: asanaCreate.Run,
	}
	createCmd.Flags().StringVar(
		&asanaCreate.AuthURL,
		"auth-url",
		asanaAuthURL,
		"OAuth2 authorization URL for the Asana provider.",
	)
	createCmd.Flags().StringVar(
		&asanaCreate.TokenURL,
		"token-url",
		asanaTokenURL,
		"OAuth2 token URL for the Asana provider.",
	)
	createCmd.Flags().StringVar(
		&asanaCreate.ClientID,
		"client-id",
		"",
		"The OAuth2 client ID for the Asana provider.",
	)
	createCmd.Flags().StringVar(
		&asanaCreate.ClientSecret,
		"client-secret",
		"",
		"The OAuth2 client secret for the Asana provider.",
	)
	createCmd.Flags().StringVar(
		&asanaCreate.OrgName,
		"org",
		"",
		"The sandbox system organization name ",
	)
	createCmd.MarkFlagRequired("client-id")
	createCmd.MarkFlagRequired("client-secret")
	createCmd.MarkFlagRequired("org")

	asanaCmd := &cobra.Command{
		Use:   "asana",
		Short: "Asana login provider related commands.",
	}
	asanaCmd.AddCommand(createCmd)

	loginProviderCmd := &cobra.Command{
		Use:   "login-provider",
		Short: "Login provider related commands.",
	}
	loginProviderCmd.AddCommand(asanaCmd)

	cs := &cobra.Command{
		Use:   "cs",
		Short: "Crafting Sandbox cs commands.",
	}
	cs.AddCommand(loginProviderCmd)

	return cs
}
