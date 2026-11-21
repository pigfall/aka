package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPiArchive(t *testing.T) {
	cases := []struct {
		platform          string
		archiveName       string
		stripComponents   int
		expectUnsupported bool
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
		"":          defaultPiVersion,
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
