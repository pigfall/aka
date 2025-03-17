package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use: "aka",
	}

	cmd.AddCommand(
		k3sCmd(),
		k8sCmd(),
		nvimCmd(),
		personalizeCmd(),
		nodejsCmd(),
		cipherCmd(),
		tlstunCmd(),
		ripgrepCmd(),
		reactCmd(),
	)

	cmd.SilenceUsage = true

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
