package config

import (
	"github.com/romainframe/grunter/pkg/utils"
)

// ConfigBuilder defines the structure for functions that can match and build configurations.
type ConfigBuilder struct {
	// Matches determines if the builder should process the given configuration.
	Matches func(Config) bool
	// Build processes and potentially modifies the configuration.
	Build func(Config) (Config, error)
}

// builders holds all registered ConfigBuilders.
// Each builder is responsible for a specific type of configuration, allowing for modular and extendable design.
var builders = []ConfigBuilder{
	K8sConfigBuilder, // Kubernetes-specific configuration builder.
}

// Build iterates over all registered builders, applying those that match the current configuration.
// It ensures that the configuration meets all requirements before attempting to build.
func (c Config) Build(extraBuilders ...ConfigBuilder) (Config, error) {
	// Template is required for building the config; return an error if it's missing.
	if c.Template == "" {
		return c, ErrTemplateRequired
	}

	// Initialize Locals map if not already done. This avoids nil map assignments.
	if c.Locals == nil {
		c.Locals = make(map[string]string)
	}

	// Iterate over each registered builder. If a builder matches the current configuration,
	// it attempts to build the configuration.
	for _, builder := range builders {
		if builder.Matches(c) {
			var err error
			c, err = builder.Build(c)
			if err != nil {
				// Return immediately if any builder encounters an error.
				return c, utils.WrapError(ErrBuildFailed, err)
			}
		}
	}

	// Iterate over any extra builders provided to the function.
	for _, builder := range extraBuilders {
		if builder.Matches(c) {
			var err error
			c, err = builder.Build(c)
			if err != nil {
				// Return immediately if any builder encounters an error.
				return c, utils.WrapError(ErrBuildFailed, err)
			}
		}
	}

	// Return the potentially modified configuration and nil error if the process completes successfully.
	return c, nil
}
