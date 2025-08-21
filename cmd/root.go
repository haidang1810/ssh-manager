package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sm/internal/config"
	"sm/internal/ssh"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sm [command] [args...]",
	Short: "A command-line tool to manage SSH connections efficiently.",
	Long: `SM is a CLI tool that helps you save, organize, and 
quickly connect to your remote servers via SSH.`, // Corrected: Removed unnecessary escaping for newline
	Args: cobra.ArbitraryArgs, // Allow any arguments
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			connName := args[0]

			// Check if the argument is a known subcommand
			for _, subCmd := range cmd.Commands() {
				if subCmd.Name() == connName {
					// It's a subcommand, let Cobra handle it normally
					return cmd.Help() // This will show help for the subcommand
				}
			}

			// It's not a known subcommand, so try to treat it as a connection name
			cfg, err := config.GetConfig()
			if err != nil {
				return fmt.Errorf("failed to get config: %w", err)
			}

			conn, exists := cfg.Connections[connName]
			if exists {
				// This is the shorthand. Execute the connect command logic directly.
				fmt.Printf("Connecting to %s (%s@%s)... (shorthand)\n", conn.Name, conn.User, conn.Host)
				if err := ssh.Connect(&conn); err != nil {
					return fmt.Errorf("ssh connection failed: %w", err)
				}
				fmt.Println("Connection closed.")
				return nil
			}
		}
		// If no connection name is provided or found, show help
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports a global flag that will be valid for all
	// subcommands, e.g:
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sm/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".sm" (without extension).
		viper.AddConfigPath(fmt.Sprintf("%s/.sm", home))
		viper.AddConfigPath(fmt.Sprintf("%s/.config/sm", home))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("SM")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
