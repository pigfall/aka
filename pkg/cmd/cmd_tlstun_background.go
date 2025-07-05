//go:build !windows

package cmd

import (
	"os/exec"
	"syscall"
)

func setRunInBackground(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
}
