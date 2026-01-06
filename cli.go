package main

import (
	"github.com/spf13/pflag"
)

type CLIOptions struct {
	ShowVersion      bool
	ConfigPath       string
	WatchFile        string
	ClipboardBackend string
}

func ParseCLI(args []string) (*CLIOptions, error) {
	fs := pflag.NewFlagSet("clipboard-txt-watcher", pflag.ContinueOnError)
	opts := &CLIOptions{}

	fs.BoolVarP(&opts.ShowVersion, "version", "v", false, "show version")
	fs.StringVarP(&opts.ConfigPath, "config", "c", "", "path to config file")
	fs.StringVarP(&opts.WatchFile, "file", "f", "", "path to file to watch")
	fs.StringVarP(&opts.ClipboardBackend, "backend", "b", "", "clipboard backend (wayland or x11)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	return opts, nil
}
