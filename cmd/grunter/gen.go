package main

import (
	"fmt"

	"github.com/romainframe/grunter/pkg/cmds"
	"github.com/romainframe/grunter/pkg/utils"
	"github.com/spf13/cobra"
)

// Predefined errors for file operations.
var (
	// ErrGenConfig is returned when Terragrunt configuration generation fails.
	ErrGenConfig = fmt.Errorf("⛔️ command 'gen' failed")
)

// genCmd represents the grunt command
var genCmd = &cobra.Command{
	Use: "gen",
	// Version: grunter.Version,
	Short: "Generate Terragrunt configuration files",
	Long: `Generate Terragrunt configuration files based on the provided config input.

This command processes a JSON or YAML file containing the necessary configuration
information and generates a corresponding Terragrunt configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Extract flag values
		inputPath, _ := cmd.Flags().GetString("input")
		outputPath, _ := cmd.Flags().GetString("output")

		err := cmds.Gen(inputPath, outputPath)
		if err != nil {
			return utils.WrapError(ErrGenConfig, err)
		}

		fmt.Printf("🎉 Terragrunt configuration successfully generated at '%s'\n", outputPath)
		return nil
	},
}

func init() {
	// Assuming infraCmd is the parent command to which genCmd is added
	rootCmd.AddCommand(genCmd)

	// Here we define the flags for genCmd
	genCmd.Flags().StringP("input", "i", "", "Path to the input configuration file")
	genCmd.Flags().StringP("output", "o", "terragrunt.hcl", "Path for the output Terragrunt configuration file (default is current directory)")
}
