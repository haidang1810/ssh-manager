package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"ssh-manager/internal/config"
)

// keysListCmd represents the list command for keys
var keysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all managed SSH keys",
	Long:  `Lists all SSH keys managed by ssh-manager, including their names, paths, and types.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		if len(cfg.SSHKeys) == 0 {
			fmt.Println("No SSH keys found. Use 'ssh-manager keys add' or 'ssh-manager keys generate' to create one.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tPATH\tTYPE")
		for _, key := range cfg.SSHKeys {
			fmt.Fprintf(w, "%s\t%s\t%s\n", key.Name, key.Path, key.Type)
		}
		w.Flush()

		return nil
	},
}

func init() {
	keysCmd.AddCommand(keysListCmd)
}
