package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPiArchive(t *testing.T) {
	cases := []struct {
		platform        string
		archiveName     string
		stripComponents int
	}{
		{platform: "linux-amd64", archiveName: "pi-linux-x64.tar.gz", stripComponents: 1},
		{platform: "linux-arm64", archiveName: "pi-linux-arm64.tar.gz", stripComponents: 1},
		{platform: "darwin-amd64", archiveName: "pi-darwin-x64.tar.gz", stripComponents: 1},
		{platform: "darwin-arm64", archiveName: "pi-darwin-arm64.tar.gz", stripComponents: 1},
		{platform: "windows-amd64", archiveName: "pi-windows-x64.zip", stripComponents: 0},
		{platform: "windows-arm64", archiveName: "pi-windows-arm64.zip", stripComponents: 0},
	}

	for _, tc := range cases {
		archiveName, stripComponents, err := piArchive(tc.platform)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tc.platform, err)
			continue
		}
		if archiveName != tc.archiveName {
			t.Errorf("%s: archive name = %q, want %q", tc.platform, archiveName, tc.archiveName)
		}
		if stripComponents != tc.stripComponents {
			t.Errorf("%s: strip components = %d, want %d", tc.platform, stripComponents, tc.stripComponents)
		}
	}

	if _, _, err := piArchive("freebsd-amd64"); err == nil {
		t.Fatal("expected error for unsupported platform, got nil")
	}
}

func TestNormalizePiVersion(t *testing.T) {
	cases := map[string]string{
		"":          "",
		"1.6.0":     "1.6.0",
		"v1.6.0":    "1.6.0",
		"  v1.6.0 ": "1.6.0",
	}

	for in, want := range cases {
		if got := normalizePiVersion(in); got != want {
			t.Errorf("normalizePiVersion(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLatestPiVersionFollowsRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, piReleasesBaseURL+"/tag/v1.6.0", http.StatusFound)
	}))
	defer server.Close()

	original := piLatestReleaseURL
	piLatestReleaseURL = server.URL
	defer func() { piLatestReleaseURL = original }()

	version, err := latestPiVersion()
	if err != nil {
		t.Fatalf("latestPiVersion() returned error: %v", err)
	}
	if version != "1.6.0" {
		t.Fatalf("latestPiVersion() = %q, want %q", version, "1.6.0")
	}
}

func TestLatestPiVersionFailsWithoutRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html>no releases</html>"))
	}))
	defer server.Close()

	original := piLatestReleaseURL
	piLatestReleaseURL = server.URL
	defer func() { piLatestReleaseURL = original }()

	if _, err := latestPiVersion(); err == nil {
		t.Fatal("expected error when the endpoint does not redirect, got nil")
	}
}

func TestInstalledPiVersion(t *testing.T) {
	home := t.TempDir()
	piDir := filepath.Join(home, "tools", "pi")

	if got := installedPiVersion(piDir); got != "" {
		t.Fatalf("installedPiVersion() for missing package.json = %q, want empty", got)
	}

	if err := os.MkdirAll(piDir, 0755); err != nil {
		t.Fatal(err)
	}
	pkg := []byte(`{"name":"@earendil-works/pi-coding-agent","version":"1.4.2"}`)
	if err := os.WriteFile(filepath.Join(piDir, "package.json"), pkg, 0644); err != nil {
		t.Fatal(err)
	}

	if got := installedPiVersion(piDir); got != "1.4.2" {
		t.Fatalf("installedPiVersion() = %q, want %q", got, "1.4.2")
	}
}

func TestAppendPathToUserShellRcIsIdempotent(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	bashrc := filepath.Join(home, ".bashrc")
	if err := os.WriteFile(bashrc, []byte("# user config\n"), 0644); err != nil {
		t.Fatal(err)
	}

	binDir := filepath.Join(home, "tools", "pi")
	for i := 0; i < 2; i++ {
		if err := appendPathToUserShellRc(binDir); err != nil {
			t.Fatalf("appendPathToUserShellRc() returned error: %v", err)
		}
	}

	content, err := os.ReadFile(bashrc)
	if err != nil {
		t.Fatal(err)
	}
	exportLine := "export PATH=" + binDir + ":$PATH"
	if count := strings.Count(string(content), exportLine); count != 1 {
		t.Fatalf("export line appears %d times, want 1; content:\n%s", count, content)
	}
	if !strings.Contains(string(content), "# user config") {
		t.Fatal("existing rc content was lost")
	}
}
