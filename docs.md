# SM - Project Documentation

## Project Overview

**Project Name:** SM  
**Language:** Go (Golang)  
**Type:** CLI Terminal Tool  
**Target:** Cross-platform (Linux, macOS, Windows)  
**Distribution:** Single binary executable  

### Core Purpose
A command-line tool to manage SSH connections efficiently, allowing users to save, organize, and quickly connect to remote servers without remembering complex SSH commands and configurations.

## Project Structure

```
sm/
├── cmd/
│   ├── root.go              # Root command setup
│   ├── add.go               # Add new SSH connection
│   ├── list.go              # List all connections
│   ├── connect.go           # Connect to saved server
│   ├── remove.go            # Remove connection
│   ├── edit.go              # Edit existing connection
│   └── export.go            # Export/import configurations
├── internal/
│   ├── config/
│   │   ├── config.go        # Configuration management
│   │   └── storage.go       # Data persistence layer
│   ├── ssh/
│   │   ├── connection.go    # SSH connection logic
│   │   ├── client.go        # SSH client wrapper
│   │   └── keys.go          # SSH key management
│   ├── models/
│   │   └── connection.go    # Data models
│   └── utils/
│       ├── crypto.go        # Encryption for sensitive data
│       ├── validation.go    # Input validation
│       └── terminal.go      # Terminal UI helpers
├── pkg/
│   └── version/
│       └── version.go       # Version information
├── configs/
│   └── example.yaml         # Example configuration
├── scripts/
│   ├── build.sh            # Build script for all platforms
│   ├── install.sh          # Installation script
│   └── release.sh          # Release automation
├── docs/
│   ├── USAGE.md            # Usage documentation
│   ├── CONFIGURATION.md    # Configuration guide
│   └── DEVELOPMENT.md      # Development setup
├── .github/
│   └── workflows/
│       └── release.yml     # GitHub Actions for releases
├── go.mod
├── go.sum
├── main.go                 # Application entry point
├── Makefile               # Build automation
├── README.md              # Project README
└── LICENSE                # MIT License
```

## Technical Requirements

### Dependencies
- **CLI Framework:** `github.com/spf13/cobra` - Modern CLI framework
- **Configuration:** `github.com/spf13/viper` - Configuration management
- **SSH Client:** `golang.org/x/crypto/ssh` - Official SSH implementation
- **Encryption:** `golang.org/x/crypto` - For secure credential storage
- **Terminal UI:** `github.com/manifoldco/promptui` - Interactive prompts
- **File Operations:** Standard library `os`, `filepath`
- **JSON/YAML:** `gopkg.in/yaml.v3` for configuration files

### Go Version
- **Minimum:** Go 1.19+
- **Recommended:** Go 1.21+

## Core Features Specification

### 1. Connection Management
```bash
# Add new connection
sm add <name> --host <host> --user <user> [--port <port>] [--key <path>] [--pass <password>]

# List all connections
sm list [--format table|json]

# Connect to server
sm connect <name>
sm <name>  # shorthand

# Remove connection
sm remove <name>

# Edit existing connection
sm edit <name>
```

### 2. Configuration Management
```bash
# Show configuration
sm config show

# Set global defaults
sm config set default-user myuser
sm config set default-port 2222

# Export/Import
sm export --output backup.yaml
sm import --input backup.yaml
```

### 3. SSH Key Management
```bash
# List SSH keys
sm keys list

# Add SSH key
sm keys add --name work --path ~/.ssh/work_rsa

# Generate new key pair
sm keys generate --name newkey --type rsa --bits 4096
```

## Data Models

### Connection Model
```go
type Connection struct {
    Name        string            `json:"name" yaml:"name"`
    Host        string            `json:"host" yaml:"host"`
    Port        int               `json:"port" yaml:"port"`
    User        string            `json:"user" yaml:"user"`
    KeyPath     string            `json:"key_path,omitempty" yaml:"key_path,omitempty"`
    Password    string            `json:"password,omitempty" yaml:"password,omitempty"` // encrypted
    Tags        []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
    Description string            `json:"description,omitempty" yaml:"description,omitempty"`
    LastUsed    time.Time         `json:"last_used" yaml:"last_used"`
    CreatedAt   time.Time         `json:"created_at" yaml:"created_at"`
    Extra       map[string]string `json:"extra,omitempty" yaml:"extra,omitempty"`
}
```

### Configuration Model
```go
type Config struct {
    DefaultUser     string                `yaml:"default_user"`
    DefaultPort     int                   `yaml:"default_port"`
    DefaultKeyPath  string                `yaml:"default_key_path"`
    ConfigPath      string                `yaml:"config_path"`
    Connections     map[string]Connection `yaml:"connections"`
    SSHKeys         map[string]SSHKey     `yaml:"ssh_keys"`
    Settings        Settings              `yaml:"settings"`
}

type Settings struct {
    EncryptPasswords bool   `yaml:"encrypt_passwords"`
    LogConnections   bool   `yaml:"log_connections"`
    LogPath         string `yaml:"log_path"`
    Editor          string `yaml:"editor"`
}
```

## Security Requirements

### 1. Password Encryption
- Use AES-256-GCM for password encryption
- Store encryption key in system keyring (Linux: libsecret, macOS: Keychain, Windows: Credential Manager)
- Never store plaintext passwords

### 2. SSH Key Security
- Support for RSA, ECDSA, Ed25519 keys
- Validate key permissions (600/400)
- Warn about world-readable keys

### 3. Configuration Security
- Config file permissions: 600 (owner read/write only)
- Secure temporary file creation
- Clear sensitive data from memory after use

## CLI Design Principles

### 1. Usability
- Intuitive command structure
- Interactive prompts for missing information
- Helpful error messages with suggestions
- Tab completion support

### 2. Output Formats
- Default: Human-readable table format
- JSON: For scripting and integration
- YAML: For configuration export
- Colorized output (disable with `--no-color`)

### 3. Configuration
- Config file locations (in priority order):
  1. `--config` flag
  2. `$SM_CONFIG`
  3. `$HOME/.sm/config.yaml`
  4. `$HOME/.config/sm/config.yaml`

## Build and Release Specification

### Build Targets
```makefile
# Makefile targets
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

build-all:
	@for platform in $(PLATFORMS); do 
		GOOS=${platform%/*} GOARCH=${platform#*/} 
		go build -ldflags "-X main.version=$(VERSION)" 
		-o bin/sm-${platform%/*}-${platform#*/} .;
		done
```

### Release Strategy
1. **GitHub Releases:** Automated via GitHub Actions
2. **Package Managers:**
   - Homebrew formula for macOS
   - AUR package for Arch Linux
   - .deb package for Ubuntu/Debian
3. **Direct Download:** Pre-built binaries for all platforms

## Development Guidelines

### Code Style
- Follow Go conventions (`gofmt`, `golint`)
- Use meaningful variable names
- Write comprehensive tests (target: >80% coverage)
- Document all exported functions

### Error Handling
- Use wrapped errors with context
- Provide actionable error messages
- Log errors appropriately
- Graceful degradation where possible

### Testing Strategy
- Unit tests for all core logic
- Integration tests for SSH connections
- CLI tests using golden files
- Cross-platform testing via CI

## User Experience Requirements

### 1. First-time Setup
- Interactive setup wizard: `sm `
- Import from existing SSH config
- Generate sample configuration

### 2. Discoverability
- Built-in help system
- Command suggestions for typos
- Examples in help text

### 3. Productivity Features
- Fuzzy search for connection names
- Recent connections list
- Batch operations
- Shell completions

## Performance Requirements

### 1. Startup Time
- Cold start: < 100ms
- Binary size: < 25MB
- Memory usage: < 10MB at startup

### 2. Connection Speed
- Connection establishment: < 2s for local networks
- Configuration loading: < 50ms
- List operations: < 100ms for 1000+ connections

## Compatibility Requirements

### Operating Systems
- Linux: Ubuntu 18.04+, RHEL 7+, Arch Linux
- macOS: 10.15+ (Catalina and newer)
- Windows: Windows 10+

### SSH Compatibility
- SSH protocol version 2.0
- Support common key exchange algorithms
- Compatible with OpenSSH server configurations

## Documentation Requirements

### User Documentation
- `README.md`: Quick start and overview
- `docs/USAGE.md`: Comprehensive usage guide
- `docs/CONFIGURATION.md`: Configuration reference
- Man page generation

### Developer Documentation
- `docs/DEVELOPMENT.md`: Setup and contribution guide
- Code comments for all public APIs
- Architecture decision records (ADRs)

## Success Metrics

### Functionality
- [ ] All core features implemented and tested
- [ ] Cross-platform compatibility verified
- [ ] Security requirements met
- [ ] Performance targets achieved

### Quality
- [ ] Test coverage > 80%
- [ ] No critical security vulnerabilities
- [ ] Documentation complete and accurate
- [ ] User feedback incorporated

### Distribution
- [ ] Automated releases working
- [ ] Package managers integration
- [ ] Installation scripts tested
- [ ] Upgrade path documented

## Implementation Phases

### Phase 1: Core Infrastructure (Week 1)
- Project structure setup
- Basic CLI framework
- Configuration management
- Data models

### Phase 2: Connection Management (Week 2)
- Add/remove/list connections
- Basic SSH connection functionality
- Configuration file I/O
- Input validation

### Phase 3: Advanced Features (Week 3)
- SSH key management
- Password encryption
- Interactive prompts
- Export/import functionality

### Phase 4: Polish and Release (Week 4)
- Comprehensive testing
- Documentation completion
- Build automation
- Release preparation

## AI Assistant Instructions

When working on this project:

1. **Follow the structure**: Implement files according to the specified project structure
2. **Implement incrementally**: Start with core functionality, then add features
3. **Test thoroughly**: Write tests for each component as you build
4. **Document as you go**: Update documentation when adding features
5. **Security first**: Always consider security implications
6. **User experience**: Think from the user's perspective for CLI design
7. **Cross-platform**: Test and consider differences between operating systems
8. **Performance**: Keep the tool fast and lightweight

When asked to implement specific components, refer to this documentation for context, requirements, and architectural decisions.
