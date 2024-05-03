package block

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Block holds the structure for application configuration, supporting nested objects
// for various configuration aspects like metadata, dependencies, and hooks.
type Block struct {
	Name         string            `json:"name"`         // Unique identifier for the block.
	Template     string            `json:"template"`     // Template path or identifier.
	Metadata     map[string]string `json:"metadata"`     // Arbitrary metadata for templating.
	Dependencies []Dependency      `json:"dependencies"` // List of external dependencies.
	Locals       map[string]string `json:"locals"`       // Local variables for templating.
	Inputs       map[string]string `json:"inputs"`       // Input variables for customization.
	BeforeHooks  []BeforeHook      `json:"beforeHooks"`  // Hooks to run before execution.
}

// BeforeHook defines a pre-execution hook with a name, commands to run, and
// scripts to execute. This can be used to prepare the environment or ensure
// prerequisites are met.
type BeforeHook struct {
	Name     string   `json:"name"`     // Unique identifier for the hook.
	Commands []string `json:"commands"` // Shell commands to run.
	Execute  []string `json:"execute"`  // Paths to scripts to execute.
}

// Dependency describes an external dependency with its source, path, and type.
// It includes an option to fetch outputs from the dependency, if applicable.
type Dependency struct {
	Name        string `json:"name"`        // Name of the dependency.
	Path        string `json:"path"`        // Location or path to the dependency.
	PathType    string `json:"pathType"`    // Type of the path (e.g., local, remote).
	WithOutputs bool   `json:"withOutputs"` // Whether to include outputs from the dependency.
}

// NewFromFile creates a Block object from a JSON or YAML file located at configPath.
// It reads the file, unmarshals into a Block struct, and processes
// it through Build() to build & validate the config.
func NewFromFile(configPath string) (Block, error) {
	// Read the entire configuration file into memory.
	fileContents, err := os.ReadFile(configPath)
	if err != nil {
		return Block{}, err // Return an empty Block and the error.
	}

	// Determine the file extension to decide on the unmarshalling method.
	ext := filepath.Ext(configPath)
	var config Block

	switch ext {
	case ".json":
		if err := json.Unmarshal(fileContents, &config); err != nil {
			return Block{}, err // Return an error if the JSON is invalid.
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(fileContents, &config); err != nil {
			return Block{}, err // Return an error if the YAML is invalid.
		}
	default:
		return Block{}, errors.New("unsupported file type")
	}

	return config, nil // Return the fully initialized Block.
}
