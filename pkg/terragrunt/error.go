package terragrunt

import "fmt"

// Predefined errors for local variable operations.
var (

	// ErrExtractLocalPath is returned when extracting a local variable path fails.
	ErrExtractLocalPath = func(part string) error {
		return fmt.Errorf("failed to extract local path from: %s", part)
	}

	// ErrValidateLocals is returned when validating local variables fails.
	ErrValidateLocals = fmt.Errorf("failed to validate local variables")
)
