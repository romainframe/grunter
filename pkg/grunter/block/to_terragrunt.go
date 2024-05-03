package block

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/romainframe/grunter/pkg/terragrunt"
	"github.com/romainframe/grunter/pkg/utils"
)

// GenTerragruntGrunt generates a Terragrunt configuration based on the Grunter's current configuration.
// It dynamically constructs the configuration by processing template sources, before hooks, dependencies,
// local variables, and inputs. The function validates necessary fields and constructs the Terragrunt configuration
// accordingly, incorporating all elements defined in the Grunter configuration into the Terragrunt configuration.
func (b Block) GenTerragruntGrunt() (terragrunt.Config, error) {
	// Create a Terragrunt configuration with initialized fields to avoid nil map/slice errors.
	tgConfig := terragrunt.Config{
		Dependencies:   []terragrunt.Dependency{},
		LocalVariables: []terragrunt.LocalVariable{},
		Inputs:         map[string]string{},
		OpenTofu: terragrunt.OpenTofu{
			BeforeHooks: []terragrunt.BeforeHook{},
		},
	}

	// Initialize a new locals search object for collecting and merging local variables.
	localsToSearch := terragrunt.NewLocalsSearch()

	// Validate the Grunter configuration's template is provided and correctly format its source.
	if b.Template == "" {
		return tgConfig, ErrTemplateRequired
	}
	tgConfig.OpenTofu.Source = formatTemplateSource(b.Template)
	localsToSearch.Add(tgConfig.OpenTofu.Source)

	// Process the Grunter configuration elements, appending them to the Terragrunt configuration.
	if err := processBeforeHooks(&tgConfig, b.BeforeHooks, localsToSearch); err != nil {
		return tgConfig, err
	}

	if err := processDependencies(&tgConfig, b.Dependencies); err != nil {
		return tgConfig, err
	}

	if err := processLocalVariables(&tgConfig, b.Locals, localsToSearch); err != nil {
		return tgConfig, err
	}

	if err := processInputs(&tgConfig, b.Inputs, localsToSearch); err != nil {
		return tgConfig, err
	}

	// Merge collected locals into the configuration's local variables.
	vars, err := localsToSearch.Merge(tgConfig.LocalVariables)
	if err != nil {
		return tgConfig, err
	}
	tgConfig.LocalVariables = vars

	return tgConfig, nil
}

// formatTemplateSource formats the source for a template, determining if it's a git source or a local template path.
func formatTemplateSource(template string) string {
	if strings.HasPrefix(template, "git::") {
		return template
	}
	return fmt.Sprintf("${local.template_root}//%s", template)
}

// processBeforeHooks processes before hooks for the Terragrunt configuration, adding them and collecting locals.
func processBeforeHooks(grunt *terragrunt.Config, beforeHooks []BeforeHook, localsSearch terragrunt.LocalsSearch) error {
	for _, bh := range beforeHooks {
		grunt.OpenTofu.BeforeHooks = append(grunt.OpenTofu.BeforeHooks, terragrunt.BeforeHook(bh))
		if err := localsSearch.Add(bh.Execute...); err != nil {
			return utils.WrapError(ErrProcessBeforeHooks(bh.Name), err)
		}
	}
	return nil
}

// processDependencies processes dependencies for the Terragrunt configuration, ensuring names are provided.
func processDependencies(grunt *terragrunt.Config, dependencies []Dependency) error {
	for _, dep := range dependencies {
		if dep.Name == "" {
			return utils.WrapError(ErrProcessDependencies(dep.Path), fmt.Errorf("dependency name is required"))
		}
		depPath, err := transformSpecialPath(dep.Path)
		if err != nil {
			return utils.WrapError(ErrProcessDependencies(dep.Path), err)
		}
		grunt.Dependencies = append(grunt.Dependencies, terragrunt.Dependency{
			Name:        dep.Name,
			ConfigPath:  depPath,
			SkipOutputs: !dep.WithOutputs,
		})
	}
	return nil
}

// processLocalVariables processes local variables for the Terragrunt configuration, adding them and collecting locals.
func processLocalVariables(grunt *terragrunt.Config, locals map[string]string, localsSearch terragrunt.LocalsSearch) error {
	for localKey, localValue := range locals {
		grunt.LocalVariables = append(grunt.LocalVariables, terragrunt.LocalVariable{
			Name:  localKey,
			Value: localValue,
		})
		if err := localsSearch.Add(localValue); err != nil {
			return utils.WrapError(ErrProcessLocals(localKey, localValue), err)
		}
	}
	return nil
}

// processInputs processes inputs for the Terragrunt configuration, adding them and collecting locals.
func processInputs(grunt *terragrunt.Config, inputs map[string]string, localsSearch terragrunt.LocalsSearch) error {
	for inKey, inValue := range inputs {
		v := inValue
		if strings.HasPrefix(inValue, "dependency.") {
			grunt.Inputs[inKey] = inValue
			continue
		}
		if strings.HasPrefix(inValue, "values.") {
			v := "local.values.locals"
			parts := strings.Split(inValue, ".")
			for i, part := range parts {
				if i > 0 && part != " " {
					v = fmt.Sprintf("%s.%s", v, part)
				}
			}
		}
		// Local variable validation: a.b.d.c.d.e
		validLocalDefRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+(\.[a-zA-Z0-9_]+)*$`)
		if validLocalDefRegex.MatchString(inValue) {
			parts := strings.Split(inValue, ".")
			if len(parts) > 0 {
				v = fmt.Sprintf("local.%s.locals", parts[0])
				for i, part := range parts {
					if i > 0 && part != " " {
						v = fmt.Sprintf("%s.%s", v, part)
					}
				}
			}
		}
		grunt.Inputs[inKey] = v
		if err := localsSearch.Add(v); err != nil {
			return utils.WrapError(ErrProcessInput(inKey, v), err)
		}
	}
	return nil
}
