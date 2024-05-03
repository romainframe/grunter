package block

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/romainframe/grunter/pkg/utils"
)

// GruntBuilder defines the structure for functions that can match and build configurations.
type GruntBuilder struct {
	// Matches determines if the builder should process the given configuration.
	Matches func(Block) bool
	// Build processes and potentially modifies the configuration.
	Build func(Block) (Block, error)
}

// builders holds all registered GruntBuilders.
// Each builder is responsible for a specific type of configuration, allowing for modular and extendable design.
var builders = []GruntBuilder{
	// K8sGruntBuilder, // Kubernetes-specific configuration builder.
}

// Build iterates over all registered builders, applying those that match the current configuration.
// It ensures that the configuration meets all requirements before attempting to build.
func (b Block) Build(systemName string, extraBuilders ...GruntBuilder) (Block, error) {
	// Name is required for building the config; return an error if it's missing.
	if b.Name == "" {
		return b, ErrNameRequired
	}
	b.Name = normalizeName(fmt.Sprintf("%s/%s", systemName, b.Name))

	// Template is required for building the config; return an error if it's missing.
	if b.Template == "" {
		return b, ErrTemplateRequired
	}

	// Initialize Locals map if not already done. This avoids nil map assignments.
	if b.Locals == nil {
		b.Locals = make(map[string]string)
	}

	// Iterate over each registered builder. If a builder matches the current configuration,
	// it attempts to build the configuration.
	for _, builder := range builders {
		if builder.Matches(b) {
			var err error
			b, err = builder.Build(b)
			if err != nil {
				// Return immediately if any builder encounters an error.
				return b, utils.WrapError(ErrBuildFailed, err)
			}
		}
	}

	// Iterate over any extra builders provided to the function.
	for _, builder := range extraBuilders {
		if builder.Matches(b) {
			var err error
			b, err = builder.Build(b)
			if err != nil {
				// Return immediately if any builder encounters an error.
				return b, utils.WrapError(ErrBuildFailed, err)
			}
		}
	}

	// Return the potentially modified configuration and nil error if the process completes successfully.
	return b, nil
}

func normalizeName(name string) string {
	parts := strings.Split(name, "/")
	names := make([]string, 0, len(parts))
	for _, part := range parts {
		names = append(names, strcase.ToLowerCamel(part))
	}
	return strings.Join(names, "/")
}
