//go:build mage

package main

import (
	"os"
	"os/exec"
)

func init() {
	os.Mkdir("out", os.ModePerm)
}

func AKA() error {
	cmd := exec.Command("go", "build", "-x", "-v", "-o", "out/aka", "./cmd/aka/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
