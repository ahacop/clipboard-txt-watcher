default:
    @just --list

# Run tidy, lint, and tests
check:
    go mod tidy
    golangci-lint run ./...
    go test -v ./...

# Build the binary
build:
    go build -ldflags "-s -w -X main.version=$(cat VERSION)" -o clipboard-txt-watcher .

# Remove build artifacts
clean:
    rm -f clipboard-txt-watcher
    rm -f result
