package main

import (
	cmdpkg "github.com/pigfall/aka/pkg/cmd"
	"github.com/spf13/cobra"
)

func k8sCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "k8s",
	}

	exampleCmd := &cobra.Command{
		Use: "example",
	}
	exampleDeployment := cmdpkg.K8SExampleDeploymentCmd{}
	exampleDeploymentCmd := &cobra.Command{
		Use:  "deployment",
		RunE: exampleDeployment.Run,
	}

	exampleCmd.AddCommand(exampleDeploymentCmd)

	cmd.AddCommand(
		exampleCmd,
	)

	return cmd
}
