package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func downloadDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "query home dir error: %v\n", err)
		os.Exit(1)
	}
	p := filepath.Join(homeDir, ".aka/download")
	if err := os.MkdirAll(p, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir %s error: %v\n", p, err)
		os.Exit(1)
	}

	return p
}

func toolDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "query home dir error: %v", err)
		os.Exit(1)
	}
	toolPath := filepath.Join(homeDir, "tools")
	if err := os.MkdirAll(toolPath, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir %s error: %w", toolPath, err)
		os.Exit(1)
	}

	envPath := os.Getenv("PATH")
	for _, v := range strings.Split(envPath, ":") {
		if v == toolPath {
			return toolPath
		}
	}

	return toolPath
}

func download(url string, saveTo io.Writer) error {
	cmd := exec.Command(
		"curl",
		"-L",
		"--show-error",
		"--fail",
		url,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = saveTo
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("download from %s error: %v", err)
	}

	return nil
}
