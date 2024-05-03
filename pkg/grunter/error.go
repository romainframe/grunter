package grunter

import "fmt"

// Predefined errors for file operations.
var (

	// ErrCreationFailed is returned when the grunter cannot be created.
	ErrCreationFailed = fmt.Errorf("failed to create grunter")
)
