package models

import "time"

// Connection represents a single SSH connection configuration.
// It contains all necessary details to establish an SSH session.
type Connection struct {
	ID          int               `json:"id" yaml:"id"`
	Name        string            `json:"name" yaml:"name"`
	Host        string            `json:"host" yaml:"host"`
	Port        int               `json:"port" yaml:"port"`
	User        string            `json:"user" yaml:"user"`
	KeyPath     string            `json:"key_path,omitempty" yaml:"key_path,omitempty"`
	Password    string            `json:"password,omitempty" yaml:"password,omitempty"` // Should be encrypted
	Tags        []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	LastUsed    time.Time         `json:"last_used,omitempty" yaml:"last_used,omitempty"`
	    CreatedAt   int64             `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Extra       map[string]string `json:"extra,omitempty" yaml:"extra,omitempty"`
}

// SSHKey represents an SSH key managed by the tool.
// This is defined in the docs but not used in the Connection struct directly.
// It will be part of the main Config.
type SSHKey struct {
	Name string `json:"name" yaml:"name"`
	Path string `json:"path" yaml:"path"`
	Type string `json:"type" yaml:"type"` // e.g., rsa, ed25519
}

// Settings defines global application settings.
type Settings struct {
	EncryptPasswords bool   `yaml:"encrypt_passwords"`
	LogConnections   bool   `yaml:"log_connections"`
	LogPath          string `yaml:"log_path"`
	Editor           string `yaml:"editor"`
}

// Config represents the entire configuration file.
// It holds all connections, managed SSH keys, and global settings.
// Note: The name in the docs is `Config`, which can be confusing.
// I'm naming it `AppConfig` to avoid conflicts with package names.
type AppConfig struct {
	NextID         int                   `json:"next_id" yaml:"next_id"`
	DefaultUser    string                `yaml:"default_user,omitempty"`
	DefaultPort    int                   `yaml:"default_port,omitempty"`
	DefaultKeyPath string                `yaml:"default_key_path,omitempty"`
	Connections    map[string]Connection `yaml:"connections"`
	SSHKeys        map[string]SSHKey     `yaml:"ssh_keys,omitempty"`
	Settings       Settings              `yaml:"settings"`
}
