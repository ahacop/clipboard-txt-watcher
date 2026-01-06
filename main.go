package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var version = "dev"

func main() {
	opts, err := ParseCLI(os.Args[1:])
	if err != nil {
		log.Fatalf("Error parsing arguments: %v", err)
	}

	if opts.ShowVersion {
		log.Printf("clipboard-txt-watcher %s", version)
		return
	}

	// Load config from file if it exists
	var cfg *Config
	cfgPath := opts.ConfigPath
	if cfgPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to get home directory: %v", err)
		}
		cfgPath = filepath.Join(homeDir, ".config", "clipboard-txt-watcher", "config.toml")
	}

	if loadedCfg, err := LoadConfig(cfgPath); err == nil {
		cfg = loadedCfg
	} else {
		cfg = &Config{ClipboardBackend: "wayland"}
	}

	// CLI flags override config file
	if opts.WatchFile != "" {
		cfg.WatchFile = opts.WatchFile
	}
	if opts.ClipboardBackend != "" {
		cfg.ClipboardBackend = opts.ClipboardBackend
	}

	if cfg.WatchFile == "" {
		log.Fatal("No watch file specified. Use --file or config file.")
	}

	log.Printf("Watching file: %s", cfg.WatchFile)
	log.Printf("Clipboard backend: %s", cfg.ClipboardBackend)

	// Create clipboard
	cb := NewClipboard(cfg.ClipboardBackend)

	// Create watcher
	w, err := NewWatcher(cfg.WatchFile, func(content string) {
		if err := SyncToClipboard(cb, content); err != nil {
			log.Printf("Failed to sync clipboard: %v", err)
		} else {
			log.Printf("Clipboard updated from file")
		}
	})
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer func() { _ = w.Close() }()

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down...")
}
