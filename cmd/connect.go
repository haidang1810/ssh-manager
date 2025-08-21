package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"sm/internal/config"
	"sm/internal/models"
	"sm/internal/ssh"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect <name_or_id>",
	Short: "Connect to a saved SSH server",
	Long:  `Establishes an interactive SSH session with the specified server configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig() // Moved to the beginning
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		identifier := args[0]
		var conn models.Connection
		var found bool
		var connName string // To store the actual name of the found connection

		// Try to parse identifier as an ID
		id, err := strconv.Atoi(identifier)
		if err == nil { // Successfully parsed as an integer
			for name, c := range cfg.Connections {
				if c.ID == id {
					conn = c
					connName = name
					found = true
					break
				}
			}
		}

		if !found { // If not found by ID, or if identifier was not an integer, try by name
			conn, found = cfg.Connections[identifier]
			connName = identifier // If found by name, the identifier is the name
		}

		if !found {
			return errors.New("connection with this name or ID does not exist")
		}

		// Update LastUsed time
		conn.LastUsed = time.Now()
		cfg.Connections[connName] = conn // Use connName to update the map
		if err := config.SaveConfig(cfg); err != nil {
			// Log this error but don't block the connection for it
			fmt.Println("Warning: could not update last used time:", err)
		}

		fmt.Println(fmt.Sprintf("Connecting to %s (%s@%s)...", conn.Name, conn.User, conn.Host))


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
