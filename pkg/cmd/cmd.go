package cmd

import (
	"fmt"
	"os"
)

type Cmd struct{}

func (c *Cmd) FailOnError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
