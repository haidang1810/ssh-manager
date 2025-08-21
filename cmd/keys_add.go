package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"ssh-manager/internal/config"
	"ssh-manager/internal/models"
)

// keysAddCmd represents the add command for keys
var keysAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an existing SSH key",
	Long:  `Adds an existing SSH key to be managed by ssh-manager.`, // Corrected: Removed unnecessary escaping for newline
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		path, _ := cmd.Flags().GetString("path")

		if name == "" || path == "" {
			return errors.New("name and path are required")
		}

		// Check if path exists and is readable
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("key file does not exist at path: %s", path)
		}

		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		if _, exists := cfg.SSHKeys[name]; exists {
			return errors.New("SSH key with this name already exists")
		}

		newKey := models.SSHKey{
			Name: name,
			Path: path,
			Type: "unknown", // We don't parse key type here, user can edit later if needed
		}

		cfg.SSHKeys[name] = newKey

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Successfully added SSH key '%s'\n", name)
		return nil
	},
}

func init() {
	keysCmd.AddCommand(keysAddCmd)

	keysAddCmd.Flags().String("name", "", "Name for the SSH key (required)")
	keysAddCmd.Flags().String("path", "", "Path to the SSH key file (required)")

	keysAddCmd.MarkFlagRequired("name")
	keysAddCmd.MarkFlagRequired("path")
}
