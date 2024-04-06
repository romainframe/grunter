package utils

import "fmt"

func WrapError(customErr, originalErr error) error {
	return fmt.Errorf("%w: %v", customErr, originalErr)
}

// Predefined errors for file operations.
var (
	// ErrFileNotFound indicates the .hcl file was not found in the expected locations.
	ErrFileNotFound = fmt.Errorf("could not find .hcl file")

	// ErrParseFailed indicates a failure in parsing the .hcl file content.
	ErrParseFailed = fmt.Errorf("could not parse .hcl file")
)
