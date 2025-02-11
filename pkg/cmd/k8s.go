package cmd

import (
	"github.com/spf13/cobra"
)

const exampleDeploymentYamlTpl = `

`

type K8SExampleDeploymentCmd struct{}

// Implment the cobra Run interface for K8SExampleDeploymentCmd
func (c *K8SExampleDeploymentCmd) Run(cobraCmd *cobra.Command, args []string) error {
	panic("TODO")
}
