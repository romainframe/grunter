package block

import "fmt"

// Predefined errors for file operations.
var (
	// ErrGruntNotFound is returned when the configuration file is not found.
	ErrGruntNotFound = func(path string) error {
		return fmt.Errorf("config file '%s' not found", path)
	}

	// ErrNameRequired is returned when the configuration is missing a name.
	ErrNameRequired = fmt.Errorf("name is required")

	// ErrTemplateRequired is returned when the configuration is missing a template.
	ErrTemplateRequired = fmt.Errorf("template is required")

	// ErrBuildFailed is returned when the configuration building process fails.
	ErrBuildFailed = fmt.Errorf("failed to build configuration")

	// ErrMetadataRequired is returned when the configuration is missing metadata.
	ErrMetadataRequired = fmt.Errorf("metadata is required")

	// ErrMetadataKeyRequired is returned when a specific key is missing from the metadata.
	ErrMetadataKeyRequired = func(key string) error {
		return fmt.Errorf("metadata key '%s' is required", key)
	}

	// ErrKeyRequired is returned when a specific key is missing.
	ErrKeyRequired = func(key string) error {
		return fmt.Errorf("key '%s' is required", key)
	}

	// ErrFilePathNotFound is returned when the file path is not found.
	ErrFilePathNotFound = func(filepath string) error {
		return fmt.Errorf("file '%s' not found", filepath)
	}

	// ErrInvalidFile is returned when the file is invalid.
	ErrInvalidFile = func(filepath string) error {
		return fmt.Errorf("file '%s' is invalid", filepath)
	}

	// ErrBeforeHookNotFound is returned when the before hook is not found.
	ErrBeforeHookNotFound = func(key string) error {
		return fmt.Errorf("before hook for %s not found", key)
	}

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
		return fmt.Errorf("failed to process dependency with path '%s'", path)
	}
)
