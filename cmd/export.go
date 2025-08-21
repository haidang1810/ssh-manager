package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"ssh-manager/internal/config"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export all connections to a YAML file",
	Long:  `Exports all saved SSH connections and configurations to a specified YAML file or to standard output.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %w", err)
		}

		bytes, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("failed to marshal config to YAML: %w", err)
		}

		outputFile, _ := cmd.Flags().GetString("output")

		if outputFile != "" {
			err = ioutil.WriteFile(outputFile, bytes, 0644)
			if err != nil {
				return fmt.Errorf("failed to write to output file %s: %w", outputFile, err)
			}
			fmt.Printf("Successfully exported configuration to %s\n", outputFile)
		} else {
			fmt.Println(string(bytes))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringP("output", "o", "", "Output file path for the backup (default is standard output)")
}