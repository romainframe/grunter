package terragrunt

import "fmt"

// LocalFunction defines a type for functions that accept a string argument and return a string.
type LocalFunction func(string) string

var (
	// SpecialLocals is a map defining special functions used to create default locals during the LocalsSearch.Search() function
	// Each function is tailored to return specific information as a string, based on a provided argument.
	// - "values": Generates a call to read a Terragrunt configuration file, expecting the filename (without extension) as input.
	// - "email": Returns the value of the TF_VAR_EMAIL environment variable, useful for templates needing access to an email.
	// - "template_root": Provides the root directory for Terraform templates from the TF_VAR_TEMPLATE_ROOT environment variable.
	SpecialLocals = map[string]LocalFunction{
		"values": func(s string) string {
			// Generates a Terragrunt configuration file read command, formatted with the provided string.
			return fmt.Sprintf(`read_terragrunt_config("%s.hcl")`, s)
		},
		"email": func(s string) string {
			// Returns the command to get the email environment variable.
			// The argument `s` is not used here, but the function signature is kept consistent for map compatibility.
			return `get_env("TF_VAR_EMAIL", "")`
		},
		"template_root": func(s string) string {
			// Returns the command to get the template root environment variable.
			// The argument `s` is not used here, but the function signature is kept consistent for map compatibility.
			return `get_env("TF_VAR_TEMPLATE_ROOT", "")`
		},
	}
)
