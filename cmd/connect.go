package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"ssh-manager/internal/config"
	"ssh-manager/internal/ssh"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect <name>",
	Short: "Connect to a saved SSH server",
	Long:  `Establishes an interactive SSH session with the specified server configuration.`, 
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		conn, exists := cfg.Connections[name]
		if !exists {
			return errors.New("connection with this name does not exist")
		}

		// Update LastUsed time
		conn.LastUsed = time.Now()
		cfg.Connections[name] = conn
		if err := config.SaveConfig(cfg); err != nil {
			// Log this error but don't block the connection for it
			fmt.Println("Warning: could not update last used time:", err)
		}

		fmt.Printf("Connecting to %s (%s@%s)...\n", conn.Name, conn.User, conn.Host)

		// The actual connection logic is in the ssh package
		if err := ssh.Connect(&conn); err != nil {
			// The error from the ssh package is often not very user-friendly
			// on its own (e.g., "EOF"). We add context here.
			return fmt.Errorf("ssh connection failed: %w", err)
		}

		fmt.Println("Connection closed.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
