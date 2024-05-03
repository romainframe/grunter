package grunter

import (
	"fmt"
	"os"

	"github.com/romainframe/grunter/pkg/grunter/block"
	"github.com/romainframe/grunter/pkg/terragrunt"
	"github.com/romainframe/grunter/pkg/utils"
)

// Grunter encapsulates the logic for generating Terragrunt configuration files.
// It holds the path to a configuration file, a Terragrunt template, and the parsed config.
type Grunter struct {
	configPath         string
	terragruntTemplate string
	valuesTemplates    string
	Object             Object
}

// NewGrunter creates and initializes a Grunter instance.
// It verifies the existence of the config file, parses it, and prepares the Grunter.
// Returns an error if the config file doesn't exist or cannot be parsed.
func New(configPath string, extraBuilders ...block.GruntBuilder) (Grunter, error) {
	var g Grunter

	// Check if the config file exists.
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return g, utils.WrapError(block.ErrGruntNotFound(configPath), err)
	}

	// Check if file is block.yaml or stack.yaml

	// Parse the configuration file.
	obj, err := NewObjectFromFile(configPath)
	if err != nil {
		return g, fmt.Errorf("could not parse config file: %w", err)
	}
	// Process the config for any post-unmarshal setup or validation.
	obj, err = obj.Build()
	if err != nil {
		return g, err
	}

	// Initialize and return a Grunter with the parsed config.
	g = Grunter{
		configPath:         configPath,
		terragruntTemplate: terragrunt.DefaultTerragruntTemplate,
		valuesTemplates:    terragrunt.DefaultValuesTemplate,
		Object:             obj,
	}
	return g, nil
}
