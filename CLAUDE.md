# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**SM** is a CLI tool for managing SSH connections, built in Go. It allows users to store, organize, and quickly connect to remote servers with support for password encryption, SSH key management, and connection tracking.

## Build and Development Commands

### Building
```bash
# Install dependencies
go mod tidy

# Build for current platform
go build -o sm .

# Build for all platforms (Linux, macOS, Windows)
go build -o sm-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o sm-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o sm-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o sm-windows-amd64.exe .
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test ./internal/config
go test ./internal/ssh
```

## Architecture

### Package Structure

- **`main.go`**: Entry point that calls `cmd.Execute()`
- **`cmd/`**: All Cobra command definitions
  - `root.go`: Root command with shorthand connection support (allows `sm <name>` to connect directly). **Launches TUI by default when no arguments provided.**
  - Command files: `add.go`, `list.go`, `connect.go`, `edit.go`, `remove.go`, `export.go`, `import.go`
  - Key management: `keys.go`, `keys_add.go`, `keys_list.go`, `keys_generate.go`
- **`internal/config/`**: Configuration management
  - Uses Viper for config file handling
  - Config stored at `~/.ssh-manager/config.yaml` (note: different from the `~/.sm/` path mentioned in root.go)
  - Handles loading and saving YAML config
- **`internal/models/`**: Data structures
  - `Connection`: SSH connection with auto-incrementing ID, name, host, port, user, credentials, metadata
  - `AppConfig`: Main config with `NextID` counter, connections map, SSH keys, settings
- **`internal/ssh/`**: SSH functionality
  - `connection.go`: Interactive SSH session establishment with PTY support
  - Handles both password and key-based auth
  - Prompts for passphrase if private key is encrypted
  - `keys.go`: SSH key generation (RSA, Ed25519) and file writing
- **`internal/utils/`**: Utilities
  - `crypto.go`: AES-256-GCM encryption/decryption using system keyring for key storage
- **`internal/tui/`**: Terminal User Interface (Bubble Tea)
  - `tui.go`: Launch function for TUI
  - `model.go`: Main TUI model and state management (ViewState, Model)
  - `connection_list.go`: Connection list view with interactive table
  - `connection_form.go`: Form for adding/editing connections
  - `dialog.go`: Delete confirmation dialog
  - `styles.go`: Lipgloss styling definitions

### Key Design Patterns

**Auto-incrementing IDs**: Connections have both a name (map key) and numeric ID. The `AppConfig.NextID` field tracks the next available ID. Commands accept either name or ID as identifiers.

**Dual authentication support**: Connections can use either SSH keys (`KeyPath`) or passwords (`Password`). Passwords are encrypted using AES-256-GCM with keys stored in the system keyring (via `99designs/keyring`).

**Shorthand command**: The root command's `RunE` checks if a single argument matches a connection name. If so, it connects directly without requiring `sm connect <name>` - just `sm <name>` works.

**Interactive prompts**: Missing required fields trigger interactive prompts using `promptui` rather than failing validation.

**TUI (Terminal User Interface)**: Built with Bubble Tea, the TUI provides an interactive interface for managing connections:
- Launched by default when running `sm` or `sm list` without arguments
- Features: arrow key navigation, connection listing with table, add/edit/delete operations
- View states managed by `ViewState` enum: `ViewList`, `ViewAddForm`, `ViewEditForm`, `ViewDeleteConfirm`
- Sub-models: `listModel` (table view), `formModel` (add/edit forms), `dialogModel` (delete confirmation)
- When user selects "Connect" in TUI, the connection is saved in `Model.connectOnExit` and executed after TUI exits
- CLI commands (like `sm add`, `sm connect <name>`) still work alongside TUI

**Config path discrepancy**: `cmd/root.go` searches for config in `~/.sm/` and `~/.config/sm/`, but `internal/config/config.go` uses `~/.ssh-manager/`. This is a known inconsistency in the codebase.

## Important Implementation Details

### Connection Identification
When implementing features that reference connections, support both name and ID lookup:
```go
// Try parsing as ID first
id, err := strconv.Atoi(identifier)
if err == nil {
    for name, c := range cfg.Connections {
        if c.ID == id {
            conn = c
            connName = name
            found = true
            break
        }
    }
}
// Fall back to name lookup
if !found {
    conn, found = cfg.Connections[identifier]
}
```

### Password Encryption
Passwords are encrypted when added/edited and decrypted on connect. Handle decryption failures gracefully (backward compatibility with plaintext):
```go
decrypted, err := utils.Decrypt(conn.Password)
if err != nil {
    // Assume plaintext for backward compatibility
    decrypted = conn.Password
}
```

### SSH Key Passphrase Handling
The `ssh.Connect()` function tries parsing keys without passphrase first, then prompts if needed. This is in `internal/ssh/connection.go:42-55`.

### Security Note
The codebase uses `ssh.InsecureIgnoreHostKey()` for convenience. This is documented as a security risk in `internal/ssh/connection.go:72`.
