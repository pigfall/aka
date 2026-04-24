package cmd

import (
    "archive/tar"
    "archive/zip"
    "compress/gzip"
    "fmt"
    "io"
    "math"
    "net"
    "net/http"
    "net/url"
    "os"
    "path"
    "path/filepath"
    "strings"
    "time"

    xzpkg "github.com/ulikunitz/xz"
    "golang.org/x/net/proxy"
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
	client, err := httpClientFromEnv()
	if err != nil {
		return fmt.Errorf("creating http client: %w", err)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request for %s error: %w", url, err)
	}
	resp, err := client.Do(req)
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

// httpClientFromEnv returns an *http.Client that routes requests through a
// SOCKS5 proxy when one is configured via environment variables. It supports
// ALL_PROXY, all_proxy, SOCKS5_PROXY and socks5_proxy. If no proxy env var is
// set, the default http.Client is returned.
func httpClientFromEnv() (*http.Client, error) {
	envVars := []string{"ALL_PROXY", "all_proxy", "SOCKS5_PROXY", "socks5_proxy"}
	var proxyVal string
	for _, k := range envVars {
		if v := os.Getenv(k); v != "" {
			proxyVal = v
			break
		}
	}
	if proxyVal == "" {
		return http.DefaultClient, nil
	}

	// Accept both bare host:port and full URLs like socks5://host:port
	u, err := url.Parse(proxyVal)
	if err != nil || u.Scheme == "" && !strings.Contains(proxyVal, ":") {
		// treat as plain host:port
		u = &url.URL{Host: proxyVal, Scheme: "socks5"}
	}

	// Only support socks5 proxies here.
	if u.Scheme != "socks5" && u.Scheme != "socks5h" {
		// For non-socks5 proxies, fall back to default client.
		return http.DefaultClient, nil
	}

	// Extract auth if present
	var auth *proxy.Auth
	if u.User != nil {
		pw, _ := u.User.Password()
		auth = &proxy.Auth{User: u.User.Username(), Password: pw}
	}

	// Ensure host includes port; default to 1080 when absent.
	host := u.Host
	if !strings.Contains(host, ":") {
		host = net.JoinHostPort(host, "1080")
	}

	dialer, err := proxy.SOCKS5("tcp", host, auth, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("create socks5 dialer: %w", err)
	}

	// net/http.Transport prefers DialContext, but the SOCKS5 dialer provides
	// Dial. Using Transport.Dial is acceptable here for simplicity.
	tr := &http.Transport{
		Dial:                dialer.Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &http.Client{Transport: tr}, nil
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

// UnpackArchive extracts common archive formats (.tar, .tar.gz, .tar.xz, .zip)
// into dstDir. stripComponents behaves like tar's --strip-components: it
// removes the specified number of leading path elements from each entry.
func UnpackArchive(srcPath, dstDir string, stripComponents int) error {
    lower := strings.ToLower(srcPath)
    switch {
    case strings.HasSuffix(lower, ".tar.gz") || strings.HasSuffix(lower, ".tgz"):
        return unpackTarGz(srcPath, dstDir, stripComponents)
    case strings.HasSuffix(lower, ".tar.xz"):
        return unpackTarXz(srcPath, dstDir, stripComponents)
    case strings.HasSuffix(lower, ".tar"):
        return unpackTar(srcPath, dstDir, stripComponents)
    case strings.HasSuffix(lower, ".zip"):
        return unpackZip(srcPath, dstDir, stripComponents)
    default:
        return fmt.Errorf("unsupported archive format: %s", srcPath)
    }
}

func stripPathElements(p string, strip int) string {
    // tar/zip entries always use forward slashes
    p = path.Clean(p)
    if p == "." || p == "/" {
        return ""
    }
    parts := strings.Split(p, "/")
    if strip >= len(parts) {
        return ""
    }
    return filepath.FromSlash(strings.Join(parts[strip:], "/"))
}

func unpackTarGz(srcPath, dstDir string, strip int) error {
    f, err := os.Open(srcPath)
    if err != nil {
        return err
    }
    defer f.Close()
    gr, err := gzip.NewReader(f)
    if err != nil {
        return err
    }
    defer gr.Close()
    tr := tar.NewReader(gr)
    return unpackTarReader(tr, dstDir, strip)
}

func unpackTarXz(srcPath, dstDir string, strip int) error {
    f, err := os.Open(srcPath)
    if err != nil {
        return err
    }
    defer f.Close()
    xr, err := xzpkg.NewReader(f)
    if err != nil {
        return err
    }
    tr := tar.NewReader(xr)
    return unpackTarReader(tr, dstDir, strip)
}

func unpackTar(srcPath, dstDir string, strip int) error {
    f, err := os.Open(srcPath)
    if err != nil {
        return err
    }
    defer f.Close()
    tr := tar.NewReader(f)
    return unpackTarReader(tr, dstDir, strip)
}

func unpackTarReader(tr *tar.Reader, dstDir string, strip int) error {
    for {
        hdr, err := tr.Next()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        name := stripPathElements(hdr.Name, strip)
        if name == "" {
            // stripped away entire path
            continue
        }
        target := filepath.Join(dstDir, name)
        switch hdr.Typeflag {
        case tar.TypeDir:
            if err := os.MkdirAll(target, os.FileMode(hdr.Mode)); err != nil {
                return err
            }
        case tar.TypeReg, tar.TypeRegA:
            if err := os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
                return err
            }
            out, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(hdr.Mode))
            if err != nil {
                return err
            }
            if _, err := io.Copy(out, tr); err != nil {
                out.Close()
                return err
            }
            out.Close()
        case tar.TypeSymlink:
            // attempt to create symlink; ignore errors on platforms that
            // don't support it.
            if err := os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
                return err
            }
            _ = os.Remove(target)
            if err := os.Symlink(hdr.Linkname, target); err != nil {
                // ignore symlink creation errors
            }
        default:
            // ignore other types
        }
    }
    return nil
}

func unpackZip(srcPath, dstDir string, strip int) error {
    zr, err := zip.OpenReader(srcPath)
    if err != nil {
        return err
    }
    defer zr.Close()
    for _, f := range zr.File {
        name := stripPathElements(f.Name, strip)
        if name == "" {
            continue
        }
        target := filepath.Join(dstDir, name)
        if f.FileInfo().IsDir() {
            if err := os.MkdirAll(target, f.Mode()); err != nil {
                return err
            }
            continue
        }
        if err := os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
            return err
        }
        in, err := f.Open()
        if err != nil {
            return err
        }
        out, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, f.Mode())
        if err != nil {
            in.Close()
            return err
        }
        if _, err := io.Copy(out, in); err != nil {
            in.Close()
            out.Close()
            return err
        }
        in.Close()
        out.Close()
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
