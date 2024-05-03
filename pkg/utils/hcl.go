package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type HCL struct {
	KeyValues map[string]interface{}
}

// NewHCL creates a new instance of HCL with the provided key-values.
func NewHCL(keyValues map[string]interface{}) HCL {
	return HCL{KeyValues: keyValues}
}

// ParseHCL attempts to parse an HCL file at the given path into an HCL instance.
func ParseHCL(path string) (HCL, error) {
	parser := hclparse.NewParser()
	file, diag := parser.ParseHCLFile(path)
	if diag.HasErrors() {
		return HCL{}, fmt.Errorf("failed to parse HCL file: %s", diag.Error())
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return HCL{}, errors.New("expected *hclsyntax.Body type")
	}

	return NewHCL(getBody(body)), nil
}

func getBody(body *hclsyntax.Body) map[string]interface{} {
	result := make(map[string]interface{})
	for _, block := range body.Blocks {
		result[block.Type] = getBody(block.Body)
	}
	for _, attribute := range body.Attributes {
		result[attribute.Name] = getAttribute(attribute)
	}
	return result
}

func getAttribute(attribute *hclsyntax.Attribute) interface{} {
	val, _ := attribute.Expr.Value(nil) // Error handling for expression evaluation can be improved if needed
	return val.AsString()
}

// Get retrieves the value associated with a given key, supporting nested keys.
func Get(h HCL, key string) (string, bool) {
	parts := strings.Split(key, ".")
	currentMap := h.KeyValues
	for i, part := range parts {
		if i == len(parts)-1 {
			value, ok := currentMap[part].(string)
			return value, ok
		} else {
			if nextMap, ok := currentMap[part].(map[string]interface{}); ok {
				currentMap = nextMap
			} else {
				return "", false
			}
		}
	}
	// This line is theoretically unreachable due to the loop logic
	return "", false
}

// GetHCLFromParent retrieves the configuration from an HCL file located within the project directory or any parent directory.
// The function searches for a file with the provided name appended with ".hcl".
// It wraps and returns any error encountered during the file search or parsing process, with additional context.
func GetHCLFromParent(name string) (HCL, error) {
	// Attempt to locate the .hcl file in the current or any parent directory.
	cloudFile, err := FindFileInParent(name+".hcl", 50)
	if err != nil {
		// Return an enhanced error message if the file is not found.
		return HCL{}, WrapError(ErrFileNotFound, err)
	}

	// Parse the found .hcl file into the HCL struct.
	h, err := ParseHCL(cloudFile)
	if err != nil {
		// Return an enhanced error message if parsing fails.
		return HCL{}, WrapError(ErrParseFailed, err)
	}

	// Return the parsed HCL data and no error on success.
	return h, nil
}
