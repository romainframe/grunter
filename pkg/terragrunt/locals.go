package terragrunt

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/romainframe/grunter/pkg/utils"
)

// LocalsSearch is a structure that holds a set of local variable names extracted from strings.
type LocalsSearch struct {
	values map[string]struct{}
}

// NewLocalsSearch initializes and returns a new instance of LocalsSearch.
func NewLocalsSearch() LocalsSearch {
	return LocalsSearch{
		values: make(map[string]struct{}),
	}
}

// Add processes and stores local variables found in the given strings.
// It supports both direct local variable names and local variables embedded within strings.
func (ls LocalsSearch) Add(values ...string) error {
	for _, value := range values {
		if err := ls.add(value); err != nil {
			return err
		}
	}
	return nil
}

// add is a helper function that extracts local variables from a single string and adds them to the set.
func (ls LocalsSearch) add(value string) error {
	value = strings.TrimSpace(value)

	// Shortcut locals
	if strings.HasPrefix(value, "values.") {
		v := "local.values.locals"
		parts := strings.Split(value, ".")
		for i, part := range parts {
			if i > 0 && part != " " {
				v = fmt.Sprintf("%s.%s", v, part)
			}
		}
		ls.values[v] = struct{}{}
		return nil
	}

	// Classic locals
	if strings.HasPrefix(value, "local.") {
		ls.values[value] = struct{}{}
		return nil
	}

	if strings.Contains(value, "local.") {
		parts := strings.Split(value, " ")
		for _, part := range parts {
			if strings.Contains(part, "${local") {
				localPath, err := extractLocalPath(part)
				if err != nil {
					return utils.WrapError(ErrExtractLocalPath(part), err)
				}
				ls.values[localPath] = struct{}{}
			}
		}
		return nil
	}

	// Local variable validation: a.b.d.c.d.e
	validLocalDefRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+(\.[a-zA-Z0-9_]+)*$`)
	if validLocalDefRegex.MatchString(value) {
		parts := strings.Split(value, ".")
		if len(parts) > 0 {
			v := fmt.Sprintf("local.%s.locals", parts[0])
			for i, part := range parts {
				if i > 0 && part != " " {
					v = fmt.Sprintf("%s.%s", v, part)
				}
			}
			ls.values[v] = struct{}{}
			return nil
		}
	}

	return nil
}

// Search retrieves the processed local variables, categorizing them as special or generic, and returns them as a map.
func (ls LocalsSearch) Search() (map[string]string, error) {
	locals := make(map[string]string)
	for local := range ls.values {
		parts := strings.Split(strings.TrimSpace(strings.ToLower(local)), ".")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid local variable: %s", local)
		}
		name := parts[1]
		if special, ok := SpecialLocals[name]; ok {
			locals[name] = special(name)
			continue
		}
		locals[name] = fmt.Sprintf(`read_terragrunt_config(find_in_parent_folders("%s.hcl"))`, name)
	}
	return locals, nil
}

// Merge integrates discovered local variables into a provided list of LocalVariable, ensuring no duplicates.
func (ls LocalsSearch) Merge(localVars []LocalVariable) ([]LocalVariable, error) {
	foundLocals, err := ls.Search()
	if err != nil {
		return nil, err
	}

	var localsToAdd []LocalVariable
	for name, value := range foundLocals {
		if !localExists(localVars, name) {
			localsToAdd = append(localsToAdd, LocalVariable{Name: name, Value: value})
		}
	}

	if len(localsToAdd) > 0 {
		localVars = append(localVars, LocalVariable{Name: "# grunted locals", Value: "begin"})
		sort.Slice(localsToAdd, func(i, j int) bool { return localsToAdd[i].Name < localsToAdd[j].Name })
		localVars = append(localVars, localsToAdd...)
		localVars = append(localVars, LocalVariable{Name: "# grunted locals", Value: "end"})
	}

	return localVars, ls.ValidateLocals(localVars)
}

// ValidateLocals ensures that all discovered local variables are present in the provided LocalVariable list.
func (ls LocalsSearch) ValidateLocals(localVars []LocalVariable) error {
	foundLocals, err := ls.Search()
	if err != nil {
		return utils.WrapError(ErrValidateLocals, err)
	}

	for name := range foundLocals {
		if !localExists(localVars, name) {
			return utils.WrapError(ErrValidateLocals, fmt.Errorf("local variable '%s' not found", name))
		}
	}

	return nil
}

// localExists checks if a given local variable name exists in the provided list of LocalVariable.
func localExists(localVars []LocalVariable, name string) bool {
	for _, localVar := range localVars {
		if localVar.Name == name {
			return true
		}
	}
	return false
}

// extractLocalPath parses a part of a string to extract the local variable path.
func extractLocalPath(part string) (string, error) {
	subParts := strings.Split(part, ".")
	localPath := "local"
	for _, subPart := range subParts[1:] {
		if trimmed := strings.Trim(strings.Split(subPart, "}")[0], " "); trimmed != "" {
			localPath = fmt.Sprintf("%s.%s", localPath, trimmed)
			break
		}
	}
	if localPath == "local" {
		return "", fmt.Errorf("invalid local variable path in: %s", part)
	}
	return localPath, nil
}
