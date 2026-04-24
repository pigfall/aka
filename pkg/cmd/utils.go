package cmd

import (
    "fmt"
    "io"
    "math"
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
		fmt.Fprintf(os.Stderr, "mkdir %s error: %v", toolPath, err)
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

    // Copy in a loop so we can show progress.
    var total int64 = resp.ContentLength
    var downloaded int64
    buf := make([]byte, 32*1024)
    for {
        n, rerr := resp.Body.Read(buf)
        if n > 0 {
            wn, werr := saveTo.Write(buf[:n])
            if werr != nil {
                return fmt.Errorf("write to target error: %w", werr)
            }
            if wn != n {
                return fmt.Errorf("short write: wrote %d of %d bytes", wn, n)
            }
            downloaded += int64(n)
            // Print progress to stderr. Use carriage return to overwrite.
            if total > 0 {
                perc := (float64(downloaded) / float64(total)) * 100
                fmt.Fprintf(os.Stderr, "\rDownloading %s / %s (%.1f%%)", humanSize(downloaded), humanSize(total), perc)
            } else {
                fmt.Fprintf(os.Stderr, "\rDownloading %s", humanSize(downloaded))
            }
        }
        if rerr != nil {
            if rerr == io.EOF {
                break
            }
            return fmt.Errorf("download from %s error: %w", url, rerr)
        }
    }
    // Newline after finished progress line.
    fmt.Fprintln(os.Stderr, "")

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

// humanSize formats bytes as a human readable string (e.g., 1.2 MB).
func humanSize(n int64) string {
    if n < 1024 {
        return fmt.Sprintf("%d B", n)
    }
    units := []string{"KB", "MB", "GB", "TB"}
    div := float64(1024)
    x := float64(n) / div
    i := 0
    for ; i < len(units)-1 && x >= 1024.0; i++ {
        x = x / 1024.0
    }
    // If number is very large, clamp to largest unit.
    if i >= len(units) {
        i = len(units) - 1
    }
    // Choose one decimal place when appropriate.
    if x < 10 {
        return fmt.Sprintf("%.1f %s", x, units[i])
    }
    return fmt.Sprintf("%.0f %s", math.Round(x), units[i])
}
