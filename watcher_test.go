package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatcher_CallsCallbackOnFileChange(t *testing.T) {
	dir := t.TempDir()
	watchFile := filepath.Join(dir, "test.txt")

	// Create initial file
	err := os.WriteFile(watchFile, []byte("initial"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	called := make(chan string, 1)
	callback := func(content string) {
		called <- content
	}

	w, err := NewWatcher(watchFile, callback)
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer func() { _ = w.Close() }()

	// Modify the file
	err = os.WriteFile(watchFile, []byte("updated"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case content := <-called:
		if content != "updated" {
			t.Errorf("expected content %q, got %q", "updated", content)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("callback was not called within timeout")
	}
}

func TestNewWatcher_ReturnsErrorForNonExistentFile(t *testing.T) {
	_, err := NewWatcher("/nonexistent/path/file.txt", func(string) {})
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestWatcher_Close_StopsWatching(t *testing.T) {
	dir := t.TempDir()
	watchFile := filepath.Join(dir, "test.txt")

	err := os.WriteFile(watchFile, []byte("initial"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	called := make(chan string, 1)
	w, err := NewWatcher(watchFile, func(content string) {
		called <- content
	})
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Write to file after close - callback should not be called
	err = os.WriteFile(watchFile, []byte("after close"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-called:
		t.Error("callback was called after Close")
	case <-time.After(100 * time.Millisecond):
		// Expected - no callback after close
	}
}
