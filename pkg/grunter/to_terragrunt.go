package grunter

import (
	"fmt"
	"strings"

	"github.com/romainframe/grunter/pkg/grunter/config"
	"github.com/romainframe/grunter/pkg/terragrunt"
	"github.com/romainframe/grunter/pkg/utils"
)

// GenTerragruntConfig generates a Terragrunt configuration based on the Grunter's current configuration.
// It dynamically constructs the configuration by processing template sources, before hooks, dependencies,
// local variables, and inputs. The function validates necessary fields and constructs the Terragrunt configuration
// accordingly, incorporating all elements defined in the Grunter configuration into the Terragrunt configuration.
func (g Grunter) GenTerragruntConfig() (terragrunt.Config, error) {
	// Create a Terragrunt configuration with initialized fields to avoid nil map/slice errors.
	tgConfig := terragrunt.Config{
		Dependencies:   []terragrunt.Dependency{},
		LocalVariables: []terragrunt.LocalVariable{},
		Inputs:         map[string]string{},
		OpenTofu: terragrunt.OFConfig{
			BeforeHooks: []terragrunt.BeforeHook{},
		},
	}

	// Initialize a new locals search object for collecting and merging local variables.
	localsToSearch := terragrunt.NewLocalsSearch()

	// Validate the Grunter configuration's template is provided and correctly format its source.
	if g.Config.Template == "" {
		return tgConfig, config.ErrTemplateRequired
	}
	tgConfig.OpenTofu.Source = formatTemplateSource(g.Config.Template)
	localsToSearch.Add(tgConfig.OpenTofu.Source)

	// Process the Grunter configuration elements, appending them to the Terragrunt configuration.
	if err := processBeforeHooks(&tgConfig, g.Config.BeforeHooks, localsToSearch); err != nil {
		return tgConfig, err
	}

	if err := processDependencies(&tgConfig, g.Config.Dependencies); err != nil {
		return tgConfig, err
	}

	if err := processLocalVariables(&tgConfig, g.Config.Locals, localsToSearch); err != nil {
		return tgConfig, err
	}

	if err := processInputs(&tgConfig, g.Config.Inputs, localsToSearch); err != nil {
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
func processBeforeHooks(config *terragrunt.Config, beforeHooks []config.BeforeHook, localsSearch terragrunt.LocalsSearch) error {
	for _, bh := range beforeHooks {
		config.OpenTofu.BeforeHooks = append(config.OpenTofu.BeforeHooks, terragrunt.BeforeHook(bh))
		if err := localsSearch.Add(bh.Execute...); err != nil {
			return utils.WrapError(ErrProcessBeforeHooks(bh.Name), err)
		}
	}
	return nil
}

// processDependencies processes dependencies for the Terragrunt configuration, ensuring names are provided.
func processDependencies(config *terragrunt.Config, dependencies []config.Dependency) error {
	for _, dep := range dependencies {
		if dep.Name == "" {
			return utils.WrapError(ErrProcessDependencies(dep.Path), fmt.Errorf("dependency name is required"))
		}
		config.Dependencies = append(config.Dependencies, terragrunt.Dependency{
			Name:        dep.Name,
			ConfigPath:  dep.Path,
			SkipOutputs: !dep.WithOutputs,
		})
	}
	return nil
}

// processLocalVariables processes local variables for the Terragrunt configuration, adding them and collecting locals.
func processLocalVariables(config *terragrunt.Config, locals map[string]string, localsSearch terragrunt.LocalsSearch) error {
	for localKey, localValue := range locals {
		config.LocalVariables = append(config.LocalVariables, terragrunt.LocalVariable{
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
func processInputs(config *terragrunt.Config, inputs map[string]string, localsSearch terragrunt.LocalsSearch) error {
	for inKey, inValue := range inputs {
		config.Inputs[inKey] = inValue
		if err := localsSearch.Add(inValue); err != nil {
			return utils.WrapError(ErrProcessInput(inKey, inValue), err)
		}
	}
	return nil
}
