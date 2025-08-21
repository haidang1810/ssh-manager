package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"ssh-manager/internal/config"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved SSH connections",
	Long:  `List all SSH connections saved in the configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		format, _ := cmd.Flags().GetString("format")

		if len(cfg.Connections) == 0 {
			fmt.Println("No connections found. Use 'ssh-manager add' to create one.")
			return nil
		}

		switch format {
		case "json":
			out, err := json.MarshalIndent(cfg.Connections, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format to json: %w", err)
			}
			fmt.Println(string(out))
		case "table":
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tUSER\tHOST\tPORT\tKEY PATH\tCREATED AT")
			for _, conn := range cfg.Connections {
				createdAtStr := "n/a"
				if conn.CreatedAt != 0 {
					createdAtStr = time.Unix(conn.CreatedAt, 0).Format(time.RFC3339)
				}
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\t%s\t%s\n", conn.ID, conn.Name, conn.User, conn.Host, conn.Port, conn.KeyPath, createdAtStr)
			}
			w.Flush()
		default:
			return fmt.Errorf("invalid format: %s. avalid formats are 'table' and 'json'", format)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("format", "f", "table", "Output format (table, json)")
}