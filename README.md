# clipboard-txt-watcher

A CLI tool that watches a text file and syncs its contents to the system clipboard whenever the file changes.

## Features

- File watching using fsnotify
- Supports Wayland (`wl-copy`/`wl-paste`), X11 (`xclip`), and macOS (`pbcopy`/`pbpaste`) clipboard backends
- Only updates clipboard when content actually changes
- Configurable via TOML config file or CLI flags

## Installation

### Using Nix Flake

Add to your NixOS configuration:

```nix
{
  inputs.clipboard-txt-watcher.url = "github:ahacop/clipboard-txt-watcher";
}
```

Then use the overlay:

```nix
{
  nixpkgs.overlays = [ inputs.clipboard-txt-watcher.overlays.default ];
  environment.systemPackages = [ pkgs.clipboard-txt-watcher ];
}
```

Or run directly:

```bash
nix run github:ahacop/clipboard-txt-watcher -- --file /path/to/file.txt
```

### Building from Source

```bash
go build -o clipboard-txt-watcher .
```

## Usage

```bash
# Watch a file using CLI flags
clipboard-txt-watcher --file /path/to/file.txt --backend wayland

# Short form
clipboard-txt-watcher -f /path/to/file.txt -b wayland

# Use a config file
clipboard-txt-watcher --config /path/to/config.toml

# Show version
clipboard-txt-watcher --version
```

### CLI Flags

| Long | Short | Description |
|------|-------|-------------|
| `--file` | `-f` | Path to the file to watch |
| `--backend` | `-b` | Clipboard backend: `wayland`, `x11`, or `darwin` |
| `--config` | `-c` | Path to config file |
| `--version` | `-v` | Show version |

## Configuration

Default config location: `~/.config/clipboard-txt-watcher/config.toml`

```toml
watch_file = "/path/to/file.txt"
clipboard_backend = "wayland"  # or "x11" or "darwin"
```

CLI flags override config file settings.

## Running as a Service (Home Manager)

The flake provides a home-manager module for running clipboard-txt-watcher as a systemd user service:

```nix
{
  inputs.clipboard-txt-watcher.url = "github:ahacop/clipboard-txt-watcher";
}
```

```nix
{
  imports = [ inputs.clipboard-txt-watcher.homeManagerModules.default ];

  services.clipboard-txt-watcher = {
    enable = true;
    watchFile = "/path/to/file.txt";
    clipboardBackend = "wayland";  # or "x11" or "darwin"
  };
}
```

This creates a systemd user service that starts automatically and restarts on failure.

## Development

```bash
# Enter development shell (provides go, golangci-lint, etc.)
nix develop

# Run all checks (tidy, lint, tests)
just check

# Build the binary
just build

# Clean build artifacts
just clean

# Run a single test
go test -v -run TestName ./...
```

## Requirements

- **Wayland**: `wl-clipboard` (provides `wl-copy` and `wl-paste`)
- **X11**: `xclip`
- **macOS**: `pbcopy` and `pbpaste` (included with macOS)

When installed via Nix on Linux, these dependencies are automatically available.

## License

GPL-3.0-or-later
