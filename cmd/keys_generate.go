package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"ssh-manager/internal/config"
	"ssh-manager/internal/models"
	"ssh-manager/internal/ssh"
)

// keysGenerateCmd represents the generate command for keys
var keysGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new SSH key pair",
	Long:  `Generates a new SSH key pair (private and public) and adds it to be managed by ssh-manager.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		keyType, _ := cmd.Flags().GetString("type")
		bits, _ := cmd.Flags().GetInt("bits")

		if name == "" {
			return errors.New("key name is required")
		}

		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		if _, exists := cfg.SSHKeys[name]; exists {
			return errors.New("SSH key with this name already exists")
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}

		sshDir := filepath.Join(homeDir, ".ssh")
		if _, err := os.Stat(sshDir); os.IsNotExist(err) {
			// Create .ssh directory if it doesn't exist
			if err := os.Mkdir(sshDir, 0700); err != nil {
				return fmt.Errorf("failed to create .ssh directory: %w", err)
			}
		}

		privateKeyPath := filepath.Join(sshDir, name)
		publicKeyPath := filepath.Join(sshDir, name+".pub")

		var privateKey interface{}
		switch keyType {
		case "rsa":
			if bits == 0 {
				bits = 2048 // Default RSA bits
			}
			privateKey, err = ssh.GenerateRSAKey(bits)
		case "ed25519":
			privateKey, err = ssh.GenerateEd25519Key()
		default:
			return fmt.Errorf("unsupported key type: %s. Supported types are 'rsa' and 'ed25519'", keyType)
		}

		if err != nil {
			return fmt.Errorf("failed to generate key: %w", err)
		}

		// Write private key
		if err := ssh.WritePrivateKey(privateKey, privateKeyPath); err != nil {
			return fmt.Errorf("failed to write private key: %w", err)
		}

		// Write public key
		if err := ssh.WritePublicKey(privateKey, publicKeyPath); err != nil {
			return fmt.Errorf("failed to write public key: %w", err)
		}

		newKey := models.SSHKey{
			Name: name,
			Path: privateKeyPath,
			Type: keyType,
		}

		cfg.SSHKeys[name] = newKey

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Successfully generated %s key '%s' at %s\n", keyType, name, privateKeyPath)
		return nil
	},
}

func init() {
	keysCmd.AddCommand(keysGenerateCmd)

	keysGenerateCmd.Flags().String("name", "", "Name for the new SSH key (required)")
	keysGenerateCmd.Flags().String("type", "rsa", "Type of key to generate (rsa or ed25519)")
	keysGenerateCmd.Flags().Int("bits", 0, "Number of bits for RSA key (e.g., 2048, 4096). Default is 2048.")

	keysGenerateCmd.MarkFlagRequired("name")
}
