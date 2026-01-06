package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	WatchFile        string `toml:"watch_file"`
	ClipboardBackend string `toml:"clipboard_backend"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := Config{
		ClipboardBackend: "wayland",
	}
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
