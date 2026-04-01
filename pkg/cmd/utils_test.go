package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownload_Success(t *testing.T) {
	expected := "hello, world!"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	}))
	defer server.Close()

	var buf bytes.Buffer
	if err := download(server.URL, &buf); err != nil {
		t.Fatalf("download() returned error: %v", err)
	}
	if buf.String() != expected {
		t.Fatalf("expected %q, got %q", expected, buf.String())
	}
}

func TestDownload_FollowsRedirects(t *testing.T) {
	expected := "redirected content"
	finalServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	}))
	defer finalServer.Close()

	redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, finalServer.URL, http.StatusFound)
	}))
	defer redirectServer.Close()

	var buf bytes.Buffer
	if err := download(redirectServer.URL, &buf); err != nil {
		t.Fatalf("download() returned error: %v", err)
	}
	if buf.String() != expected {
		t.Fatalf("expected %q, got %q", expected, buf.String())
	}
}

func TestDownload_FailsOnNon200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	var buf bytes.Buffer
	err := download(server.URL, &buf)
	if err == nil {
		t.Fatal("expected error for 404 response, got nil")
	}
}

func TestDownloadToFile_Success(t *testing.T) {
	expected := "file content here"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "testfile.txt")

	if err := downloadToFile(server.URL, filePath); err != nil {
		t.Fatalf("downloadToFile() returned error: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read downloaded file: %v", err)
	}
	if string(content) != expected {
		t.Fatalf("expected %q, got %q", expected, string(content))
	}
}

func TestDownloadToFile_CleansUpOnFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "testfile.txt")

	err := downloadToFile(server.URL, filePath)
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}

	if _, statErr := os.Stat(filePath); !os.IsNotExist(statErr) {
		t.Fatal("expected file to be cleaned up after failed download")
	}
}

func TestDownloadToFile_LargePayload(t *testing.T) {
	// Generate a 1MB payload
	payload := make([]byte, 1024*1024)
	for i := range payload {
		payload[i] = byte(i % 256)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "largefile.bin")

	if err := downloadToFile(server.URL, filePath); err != nil {
		t.Fatalf("downloadToFile() returned error: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read downloaded file: %v", err)
	}
	if !bytes.Equal(content, payload) {
		t.Fatalf("downloaded content does not match expected payload (got %d bytes, expected %d bytes)", len(content), len(payload))
	}
}
