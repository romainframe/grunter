package grunter

import "fmt"

// Predefined errors for file operations.
var (

	// ErrCreationFailed is returned when the grunter cannot be created.
	ErrCreationFailed = fmt.Errorf("failed to create grunter")

	// ErrProcessBeforeHooks is returned when processing before hooks fails.
	ErrProcessBeforeHooks = func(hook string) error {
		return fmt.Errorf("failed to process before hook '%s'", hook)
	}

	// ErrProcessLocals is returned when processing local variables fails.
	ErrProcessLocals = func(key, value string) error {
		return fmt.Errorf("failed to process local variable '%s' with value '%s'", key, value)
	}

	// ErrProcessInput is returned when processing inputs fails.
	ErrProcessInput = func(key, value string) error {
		return fmt.Errorf("failed to process input '%s' with value '%s'", key, value)
	}

	// ErrProcessDependencies is returned when processing dependencies fails.
	ErrProcessDependencies = func(path string) error {
		return fmt.Errorf("dependency path '%s' is missing a name", path)
	}
)
