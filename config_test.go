package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_ReadsWatchFile(t *testing.T) {
	// Create a temp config file
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	err := os.WriteFile(configPath, []byte(`watch_file = "/tmp/clipboard.txt"`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.WatchFile != "/tmp/clipboard.txt" {
		t.Errorf("got WatchFile=%q, want %q", cfg.WatchFile, "/tmp/clipboard.txt")
	}
}

func TestLoadConfig_ClipboardBackendDefaultsToWayland(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	err := os.WriteFile(configPath, []byte(`watch_file = "/tmp/clipboard.txt"`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.ClipboardBackend != "wayland" {
		t.Errorf("got ClipboardBackend=%q, want %q", cfg.ClipboardBackend, "wayland")
	}
}

func TestLoadConfig_ClipboardBackendCanBeOverridden(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	content := `watch_file = "/tmp/clipboard.txt"
clipboard_backend = "x11"`
	err := os.WriteFile(configPath, []byte(content), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.ClipboardBackend != "x11" {
		t.Errorf("got ClipboardBackend=%q, want %q", cfg.ClipboardBackend, "x11")
	}
}

func TestLoadConfig_ReturnsErrorForNonExistentFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.toml")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}
