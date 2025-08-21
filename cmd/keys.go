package cmd

import (
	"github.com/spf13/cobra"
)

// keysCmd represents the keys command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage SSH keys",
	Long:  `Allows adding, listing, and generating SSH keys managed by ssh-manager.`,
}

func init() {
	rootCmd.AddCommand(keysCmd)
}
