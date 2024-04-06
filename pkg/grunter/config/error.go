package config

import "fmt"

// Predefined errors for file operations.
var (
	// ErrConfigNotFound is returned when the configuration file is not found.
	ErrConfigNotFound = func(path string) error {
		return fmt.Errorf("config file '%s' not found", path)
	}

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
)
