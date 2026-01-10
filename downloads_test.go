package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetDownloadsFolder(t *testing.T) {
	// Create a temporary directory for testing
	tmpHome, err := os.MkdirTemp("", "downloads-test-get")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	// Mock UserHomeDir via environment variables
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tmpHome)

	originalUserProfile := os.Getenv("USERPROFILE")
	defer os.Setenv("USERPROFILE", originalUserProfile)
	os.Setenv("USERPROFILE", tmpHome)

	app := &App{testHomeDir: tmpHome}

	// Create a mock Downloads folder
	mockDownloads := filepath.Join(tmpHome, "Downloads")
	err = os.MkdirAll(mockDownloads, 0755)
	if err != nil {
		t.Fatalf("Failed to create mock downloads dir: %v", err)
	}

	folder, err := app.GetDownloadsFolder()
	if err != nil {
		t.Fatalf("GetDownloadsFolder failed: %v", err)
	}
	
	t.Logf("Got downloads folder: %s", folder)
	
	if folder == "" {
		t.Errorf("Expected non-empty folder path")
	}
}

func TestDownloadUpdate(t *testing.T) {
	// Create a temporary directory for testing
	tmpHome, err := os.MkdirTemp("", "download-update-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	// Mock UserHomeDir
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tmpHome)

	originalUserProfile := os.Getenv("USERPROFILE")
	defer os.Setenv("USERPROFILE", originalUserProfile)
	os.Setenv("USERPROFILE", tmpHome)

	app := NewApp()

	// Create a mock Downloads folder
	mockDownloads := filepath.Join(tmpHome, "Downloads")
	err = os.MkdirAll(mockDownloads, 0755)
	if err != nil {
		t.Fatalf("Failed to create mock downloads dir: %v", err)
	}

	// Set up a local HTTP server to serve a dummy file
	dummyContent := "this is a dummy installer"
	mux := http.NewServeMux()
	mux.HandleFunc("/AICoder-Setup.exe", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(dummyContent)))
		w.Write([]byte(dummyContent))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	// We can't easily test EventsEmit without a real Wails context,
	// but we can check if the file is downloaded correctly.
	
	destPath, err := app.DownloadUpdate(server.URL+"/AICoder-Setup.exe", "AICoder-Setup.exe")
	if err != nil {
		t.Fatalf("DownloadUpdate failed: %v", err)
	}

	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		t.Errorf("Downloaded file does not exist: %s", destPath)
	}

	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(content) != dummyContent {
		t.Errorf("Downloaded content mismatch. Got %s, expected %s", string(content), dummyContent)
	}
}

func TestCancelDownload(t *testing.T) {
	tmpHome, err := os.MkdirTemp("", "cancel-download-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	// Mock UserHomeDir
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tmpHome)

	originalUserProfile := os.Getenv("USERPROFILE")
	defer os.Setenv("USERPROFILE", originalUserProfile)
	os.Setenv("USERPROFILE", tmpHome)

	// Create a mock Downloads folder
	mockDownloads := filepath.Join(tmpHome, "Downloads")
	err = os.MkdirAll(mockDownloads, 0755)
	if err != nil {
		t.Fatalf("Failed to create mock downloads dir: %v", err)
	}

	// Slow server to allow time for cancellation
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1000000")
		w.WriteHeader(http.StatusOK)
		for i := 0; i < 100; i++ {
			w.Write(make([]byte, 1000))
			time.Sleep(10 * time.Millisecond)
		}
	}))
	defer server.Close()

	app := NewApp()
	
	fileName := "CancelTest.exe"
	
	// Start download in a goroutine
	errChan := make(chan error, 1)
	go func() {
		_, err := app.DownloadUpdate(server.URL, fileName)
		errChan <- err
	}()

	// Wait a bit and then cancel
	time.Sleep(100 * time.Millisecond)
	app.CancelDownload(fileName)

	err = <-errChan
	if err == nil || err.Error() != "download cancelled" {
		t.Errorf("Expected 'download cancelled' error, got: %v", err)
	}
}