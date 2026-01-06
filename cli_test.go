package main

import (
	"testing"
)

func TestParseCLI_VersionFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"long form", []string{"--version"}},
		{"short form", []string{"-v"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := ParseCLI(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !opts.ShowVersion {
				t.Error("expected ShowVersion to be true")
			}
		})
	}
}

func TestParseCLI_ConfigFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"long form", []string{"--config", "/path/to/config.toml"}},
		{"short form", []string{"-c", "/path/to/config.toml"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := ParseCLI(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if opts.ConfigPath != "/path/to/config.toml" {
				t.Errorf("expected ConfigPath to be '/path/to/config.toml', got '%s'", opts.ConfigPath)
			}
		})
	}
}

func TestParseCLI_FileFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"long form", []string{"--file", "/path/to/watch.txt"}},
		{"short form", []string{"-f", "/path/to/watch.txt"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := ParseCLI(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if opts.WatchFile != "/path/to/watch.txt" {
				t.Errorf("expected WatchFile to be '/path/to/watch.txt', got '%s'", opts.WatchFile)
			}
		})
	}
}

func TestParseCLI_BackendFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"long form", []string{"--backend", "x11"}},
		{"short form", []string{"-b", "x11"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := ParseCLI(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if opts.ClipboardBackend != "x11" {
				t.Errorf("expected ClipboardBackend to be 'x11', got '%s'", opts.ClipboardBackend)
			}
		})
	}
}
