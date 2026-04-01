package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download from %s error: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download from %s failed with status: %s", url, resp.Status)
	}

	if _, err := io.Copy(saveTo, resp.Body); err != nil {
		return fmt.Errorf("download from %s error: %w", url, err)
	}

	return nil
}

func downloadToFile(url string, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file %s error: %w", filePath, err)
	}
	defer f.Close()

	if err := download(url, f); err != nil {
		os.Remove(filePath)
		return err
	}

	return nil
}
