package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"sm/internal/config"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <name_or_id>",
	Short: "Remove an existing SSH connection",
	Long:  `Remove an existing SSH connection from the configuration file.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := args[0]
		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		var connName string // To store the actual name of the connection to be removed
		var found bool

		// Try to parse identifier as an ID
		id, err := strconv.Atoi(identifier)
		if err == nil { // Successfully parsed as an integer
			for name, c := range cfg.Connections {
				if c.ID == id {
					connName = name
					found = true
					break
				}
			}
		}

		if !found { // If not found by ID, or if identifier was not an integer, try by name
			if _, exists := cfg.Connections[identifier]; exists {
				connName = identifier
				found = true
			}
		}

		if !found {
			return errors.New("connection with this name or ID does not exist")
		}

		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Are you sure you want to remove connection '%s'", connName),
			IsConfirm: true,
		}

		_, err = prompt.Run()

		if err != nil {
			// User chose not to remove, or an error occurred. 
			// If the error is just that the user aborted, we don't print it.
			if err == promptui.ErrAbort {
				fmt.Println("Remove operation cancelled.")
				return nil
			}
			return fmt.Errorf("prompt failed: %w", err)
		}

		delete(cfg.Connections, connName)

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Successfully removed connection '%s'\n", connName)
		return nil
	},
}


func init() {
	rootCmd.AddCommand(removeCmd)
}