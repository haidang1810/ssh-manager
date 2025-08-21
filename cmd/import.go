package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"ssh-manager/internal/config"
	"ssh-manager/internal/models"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import connections from a YAML file",
	Long:  `Imports SSH connections from a specified YAML file, merging them with the existing configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile, _ := cmd.Flags().GetString("input")
		if inputFile == "" {
			return fmt.Errorf("input file must be specified with --input or -i")
		}

		bytes, err := ioutil.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file %s: %w", inputFile, err)
		}

		var importedCfg models.AppConfig
		if err := yaml.Unmarshal(bytes, &importedCfg); err != nil {
			return fmt.Errorf("failed to parse YAML from input file: %w", err)
		}

		currentCfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get current config: %w", err)
		}

		importedCount := 0
		skippedCount := 0
		for name, conn := range importedCfg.Connections {
			if _, exists := currentCfg.Connections[name]; exists {
				skippedCount++
				continue
			}
			currentCfg.Connections[name] = conn
			importedCount++
		}

		if err := config.SaveConfig(currentCfg); err != nil {
			return fmt.Errorf("failed to save updated config: %w", err)
		}

		fmt.Printf("Import complete. Added: %d, Skipped: %d\n", importedCount, skippedCount)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringP("input", "i", "", "Input file path for the backup (required)")
	importCmd.MarkFlagRequired("input")
}
