# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

```bash
# Run all checks (tidy, lint, tests)
just check

# Individual commands
go mod tidy
golangci-lint run ./...
go test -v ./...

# Run a single test
go test -v -run TestWatcher_CallsCallbackOnFileChange ./...

# Build the binary
go build -o clipboard-txt-watcher .

# Build with version info (as flake.nix does)
go build -ldflags "-s -w -X main.version=$(cat VERSION)" -o clipboard-txt-watcher .
```

## Architecture

This is a Go CLI tool that watches a text file and syncs its contents to the system clipboard whenever the file changes.

**Core components:**

- `main.go` - Entry point, config loading, signal handling
- `watcher.go` - File watcher using fsnotify, triggers callback on write events
- `clipboard.go` - Clipboard interface with Wayland (`wl-copy`/`wl-paste`) and X11 (`xclip`) backends
- `sync.go` - Sync logic: only writes to clipboard if content differs
- `config.go` - TOML config loading with defaults

**Testing approach:** Clipboard implementations accept injectable command executors for unit testing without real clipboard access.

## Configuration

Config file location: `~/.config/clipboard-txt-watcher/config.toml`

```toml
watch_file = "/path/to/file.txt"
clipboard_backend = "wayland"  # or "x11"
```

## NixOS/Nix

Uses a Nix flake for development environment and packaging. The flake provides an overlay for use in NixOS configurations.

## TDD Rules

Follow strict Test-Driven Development with the **Red → Green → Refactor** cycle.

**Run `just check` at every step** (tidy, lint, tests).

### The Cycle

1. **Red** - Write one test. Run `just check`. Verify it fails.
2. **Green** - Write the smallest implementation to make it pass. Run `just check`. Verify it passes.
3. **Refactor** - Clean up duplication, simplify code. Run `just check`. Verify nothing broke.

### Rules

- **One test at a time** - Do not write all tests at once. Write a single test, make it pass, then write the next.
- **Smallest implementation** - Write the minimal code needed to pass. No more.
- **`just check` at every step** - Never skip running the full check at red, green, AND refactor.
- **Actually refactor** - Don't skip the refactor step. Either refactor or explicitly decide "no refactoring needed."
- **Test error paths** - Always write tests for error handling paths, not just happy paths. If a function can return an error, test that it does so under the right conditions.
