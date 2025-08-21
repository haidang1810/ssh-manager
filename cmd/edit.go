package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"sm/internal/config"
	"sm/internal/utils"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit an existing SSH connection",
	Long:  `Edit an existing SSH connection by providing new values for the fields you want to change.`, 
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

		if cmd.Flags().Changed("host") {
			conn.Host, _ = cmd.Flags().GetString("host")
		}
		if cmd.Flags().Changed("user") {
			conn.User, _ = cmd.Flags().GetString("user")
		}
		if cmd.Flags().Changed("port") {
			conn.Port, _ = cmd.Flags().GetInt("port")
		}
		if cmd.Flags().Changed("key") {
			conn.KeyPath, _ = cmd.Flags().GetString("key")
		}
		if cmd.Flags().Changed("pass") {
			password, _ := cmd.Flags().GetString("pass")
			if password != "" {
				encryptedPass, err := utils.Encrypt(password)
				if err != nil {
					return fmt.Errorf("failed to encrypt password: %w", err)
				}
			conn.Password = encryptedPass
			} else {
			conn.Password = "" // Clear password if empty string is provided
			}
		}

		cfg.Connections[name] = conn

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Successfully updated connection '%s'\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().String("host", "", "New host name or IP address")
	editCmd.Flags().String("user", "", "New username for the connection")
	editCmd.Flags().IntP("port", "p", 0, "New port number for the connection")
	editCmd.Flags().String("key", "", "New path to the private SSH key")
	editCmd.Flags().String("pass", "", "New password for the connection")
}