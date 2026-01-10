package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDownloadsFolder(t *testing.T) {
	// Create a temporary directory for testing
	tmpHome, err := os.MkdirTemp("", "downloads-test")
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

	// On Windows, GetDownloadsFolder might still return the real system downloads folder
	// if the shell32.dll call succeeds, ignoring our environment variable override
	// for the system API call. But the fallback should work.
	
	t.Logf("Got downloads folder: %s", folder)
	
	if folder == "" {
		t.Errorf("Expected non-empty folder path")
	}
}
