package cmd

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"ssh-manager/internal/config"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove an existing SSH connection",
	Long:  `Remove an existing SSH connection from the configuration file.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		if _, exists := cfg.Connections[name]; !exists {
			return errors.New("connection with this name does not exist")
		}

		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Are you sure you want to remove connection '%s'", name),
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

		delete(cfg.Connections, name)

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Successfully removed connection '%s'\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}