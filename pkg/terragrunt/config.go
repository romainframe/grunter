package terragrunt

import "strings"

// Config stores the configuration for a Terragrunt project, including dependencies,
// local variables, OpenTofu configurations, and inputs.
type Config struct {
	Dependencies   []Dependency      `json:"dependencies"`    // List of external Terragrunt dependencies
	LocalVariables []LocalVariable   `json:"local_variables"` // Local variables specific to the Terragrunt configuration
	OpenTofu       OpenTofu          `json:"open_tofu"`       // Configuration for OpenTofu, a fictional feature or module
	Inputs         map[string]string `json:"inputs"`          // Key-value pairs for Terragrunt inputs
}

func (c Config) GetDefaultValues() map[string]string {
	defaults := make(map[string]string)
	for _, value := range c.Inputs {
		val := ""
		if strings.HasPrefix(value, "local.values.locals") {
			val = strings.TrimPrefix(value, "local.values.locals.")
		}
		if val != "" {
			defaults[val] = val
		}
	}
	return defaults
}

// Dependency defines a Terragrunt project's external dependency, including its
// name, configuration path, and whether to skip outputs.
type Dependency struct {
	Name        string `json:"name"`         // Unique identifier of the dependency
	ConfigPath  string `json:"config_path"`  // File path to the dependency's Terragrunt configuration
	SkipOutputs bool   `json:"skip_outputs"` // Indicates if outputs from this dependency should be ignored
}

// LocalVariable represents a key-value pair used as a local variable within
// a Terragrunt configuration.
type LocalVariable struct {
	Name  string `json:"name"`  // Name of the local variable
	Value string `json:"value"` // Value of the local variable
}

// BeforeHook defines a hook that runs before certain Terragrunt actions, specifying
// the hook name, commands to run, and specific commands to execute.
type BeforeHook struct {
	Name     string   `json:"name"`     // Name of the before-hook
	Commands []string `json:"commands"` // Shell commands to run as part of this hook
	Execute  []string `json:"execute"`  // Specific Terragrunt commands that trigger this hook
}

// OFConfig holds configurations related to OpenTofu within a Terragrunt project, including
// its source and any before-hooks that should run.
type OpenTofu struct {
	Source      string       `json:"source"`       // The source configuration for OpenTofu
	BeforeHooks []BeforeHook `json:"before_hooks"` // Hooks to run before certain actions are executed
}
