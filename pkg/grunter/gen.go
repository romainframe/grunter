package grunter

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/romainframe/grunter/pkg/utils"
)

// Gen generates a Terragrunt configuration file based on the Grunter's config.
// It writes the generated configuration to the specified outputPath or to
// './terragrunt.hcl' if outputPath is empty. Returns an error if the process fails.
func (g Grunter) Gen(outputPath string) ([]string, error) {
	// Use a default path if none is specified.
	if outputPath == "" {
		outputPath = "./terragrunt.hcl"
	}

	// Convert the internal config to a Terragrunt configuration.
	tgGrunts, err := g.Object.GenTerragruntGrunts(outputPath)
	if err != nil {
		return nil, fmt.Errorf("could not convert config to terragrunt config: %w", err)
	}

	generatedFiles := []string{}

	for path, tgGrunt := range tgGrunts {
		if filepath.Ext(path) == "" {
			if !utils.DoesFileOrDirExists(path) {
				err := os.MkdirAll(path, os.ModePerm)
				if err != nil {
					return nil, fmt.Errorf("could not create output directory: %w", err)
				}
			}

			valuesPath := filepath.Join(path, "values.hcl")
			if !utils.DoesFileOrDirExists(valuesPath) {
				valuesTmpl, err := template.New("values").Parse(g.valuesTemplates)
				if err != nil {
					return nil, fmt.Errorf("could not parse values template: %w", err)
				}

				// Touch a values.hcl in the folder
				valuesFile, err := os.Create(valuesPath)
				if err != nil {
					return nil, fmt.Errorf("could not create values file: %w", err)
				}
				defer func() {
					if valuesFile != nil {
						valuesFile.Close() // Ensure file closure after writing.
					}
				}()

				err = valuesTmpl.Execute(valuesFile, tgGrunt.GetDefaultValues())
				if err != nil {
					return nil, fmt.Errorf("could not execute values template: %w", err)
				}

				if valuesFile != nil {
					valuesFile.Close() // Ensure file closure after writing.
				}
			}

			// Create or overwrite the Terragrunt configuration file.
			path = filepath.Join(path, outputPath)
		}

		// Prepare the Terragrunt template.
		tmpl, err := template.New("terragrunt").Parse(g.terragruntTemplate)
		if err != nil {
			return nil, fmt.Errorf("could not parse terragrunt template: %w", err)
		}

		// Create or overwrite the Terragrunt configuration file.
		outputFile, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("could not create output file: %w", err)
		}
		defer func() {
			if outputFile != nil {
				outputFile.Close() // Ensure file closure after writing.
			}
		}()

		// Execute the template and write the output to the file.
		err = tmpl.Execute(outputFile, tgGrunt)
		if err != nil {
			return nil, fmt.Errorf("could not execute template: %w", err)
		}

		if outputFile != nil {
			outputFile.Close() // Ensure file closure after writing.
		}

		generatedFiles = append(generatedFiles, path)
	}

	return generatedFiles, nil
}
