package grunter

import (
	"fmt"
	"os"
	"text/template"

	"github.com/romainframe/grunter/pkg/grunter/config"
	"github.com/romainframe/grunter/pkg/terragrunt"
	"github.com/romainframe/grunter/pkg/utils"
)

// Grunter encapsulates the logic for generating Terragrunt configuration files.
// It holds the path to a configuration file, a Terragrunt template, and the parsed config.
type Grunter struct {
	configPath         string
	terragruntTemplate string
	Config             config.Config
}

// NewGrunter creates and initializes a Grunter instance.
// It verifies the existence of the config file, parses it, and prepares the Grunter.
// Returns an error if the config file doesn't exist or cannot be parsed.
func New(configPath string, extraBuilders ...config.ConfigBuilder) (Grunter, error) {
	var g Grunter

	// Check if the config file exists.
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return g, utils.WrapError(config.ErrConfigNotFound(configPath), err)
	}
	// Parse the configuration file.
	cfg, err := config.NewFromFile(configPath)
	if err != nil {
		return g, fmt.Errorf("could not parse config file: %w", err)
	}
	// Process the config for any post-unmarshal setup or validation.
	cfg, err = cfg.Build(extraBuilders...)
	if err != nil {
		return g, err
	}

	// Initialize and return a Grunter with the parsed config.
	g = Grunter{
		configPath:         configPath,
		terragruntTemplate: terragrunt.DefaultTerragruntTemplate,
		Config:             cfg,
	}
	return g, nil
}

// Grunt generates a Terragrunt configuration file based on the Grunter's config.
// It writes the generated configuration to the specified outputPath or to
// './terragrunt.hcl' if outputPath is empty. Returns an error if the process fails.
func (g Grunter) Grunt(outputPath string) error {
	// Use a default path if none is specified.
	if outputPath == "" {
		outputPath = "./terragrunt.hcl"
	}

	// Convert the internal config to a Terragrunt configuration.
	tgConfig, err := g.GenTerragruntConfig()
	if err != nil {
		return fmt.Errorf("could not convert config to terragrunt config: %w", err)
	}

	// Prepare the Terragrunt template.
	tmpl, err := template.New("terragrunt").Parse(g.terragruntTemplate)
	if err != nil {
		return fmt.Errorf("could not parse terragrunt template: %w", err)
	}

	// Create or overwrite the Terragrunt configuration file.
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer outputFile.Close() // Ensure file closure after writing.

	// Execute the template and write the output to the file.
	return tmpl.Execute(outputFile, tgConfig)
}
