package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"ssh-manager/internal/config"
	"ssh-manager/internal/models"
	"ssh-manager/internal/utils"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new SSH connection",
	Long:  `Add a new SSH connection to the configuration file.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		if _, exists := cfg.Connections[name]; exists {
			return errors.New("connection with this name already exists")
		}

		host, _ := cmd.Flags().GetString("host")
		user, _ := cmd.Flags().GetString("user")
		port, _ := cmd.Flags().GetInt("port")
		key, _ := cmd.Flags().GetString("key")
		password, _ := cmd.Flags().GetString("pass")

		// Interactive prompts for missing required fields
		if host == "" {
			prompt := promptui.Prompt{
				Label: "Host",
				Validate: func(input string) error {
					if len(input) == 0 {
						return errors.New("host cannot be empty")
					}
					return nil
				},
			}
			host, err = prompt.Run()
			if err != nil {
				return fmt.Errorf("prompt failed: %w", err)
			}
		}

		if user == "" {
			prompt := promptui.Prompt{
				Label: "User",
				Validate: func(input string) error {
					if len(input) == 0 {
						return errors.New("user cannot be empty")
					}
					return nil
				},
			}
			user, err = prompt.Run()
			if err != nil {
				return fmt.Errorf("prompt failed: %w", err)
			}
		}

		if port == 0 {
			prompt := promptui.Prompt{
				Label:   "Port (default: 22)",
				Default: "22",
				Validate: func(input string) error {
					if _, err := strconv.Atoi(input); err != nil {
						return errors.New("invalid port number")
					}
					return nil
				},
			}
			portStr, err := prompt.Run()
			if err != nil {
				return fmt.Errorf("prompt failed: %w", err)
			}
			port, _ = strconv.Atoi(portStr)
		}

		// Encrypt password if provided
		if password != "" {
			encryptedPass, err := utils.Encrypt(password)
			if err != nil {
				return fmt.Errorf("failed to encrypt password: %w", err)
			}
			password = encryptedPass
		}

		newConn := models.Connection{
			Name:      name,
			Host:      host,
			User:      user,
			Port:      port,
			KeyPath:   key,
			Password:  password,
			CreatedAt: time.Now().Unix(),
		}

		cfg.Connections[name] = newConn

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Successfully added connection '%s'\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().String("host", "", "Host name or IP address")
	addCmd.Flags().String("user", "", "Username for the connection")
	addCmd.Flags().IntP("port", "p", 0, "Port number for the connection (default: 22)")
	addCmd.Flags().String("key", "", "Path to the private SSH key")
	addCmd.Flags().String("pass", "", "Password for the connection (not recommended, will be stored in plaintext for now)")

	// Removed MarkFlagRequired for interactive prompts
}
