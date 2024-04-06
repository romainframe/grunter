package cmds

import (
	"fmt"

	"github.com/romainframe/grunter/pkg/grunter"
	"github.com/romainframe/grunter/pkg/utils"
)

// Predefined errors for file operations.
var (
	// ErrInitGrunter is returned when Grunter initialization fails.
	ErrInitGrunter = fmt.Errorf("⛔️ failed to initialize grunter")
	// ErrGenConfig is returned when Terragrunt configuration generation fails.
	ErrGenConfig = fmt.Errorf("⛔️ failed to generate Terragrunt configuration")
)

// Gen generates the Terragrunt configuration based on the provided input and output paths.
// If inputPath is empty, it defaults to "config.json". This function initializes Grunter
// with the given inputPath, and then calls its Grunt method to generate the configuration
// at outputPath. It handles and returns errors during the Grunter initialization and configuration
// generation process.
func Gen(inputPath, outputPath string) error {
	// Default input path to "config.json" if empty
	if inputPath == "" {
		inputPath = "config.yaml"
	}

	// Initialize Grunter with the specified inputPath
	grunter, err := grunter.New(inputPath)
	if err != nil {
		// Return an error with additional context if Grunter initialization fails
		return utils.WrapError(ErrInitGrunter, err)
	}

	// Generate the Terragrunt configuration using the initialized Grunter
	if err := grunter.Grunt(outputPath); err != nil {
		// Return an error with additional context if configuration generation fails
		return utils.WrapError(ErrGenConfig, err)
	}

	// Return nil if no errors occurred, indicating success
	return nil
}
