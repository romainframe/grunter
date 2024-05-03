package system

import (
	"fmt"

	"github.com/romainframe/grunter/pkg/terragrunt"
)

func (s System) GenTerragruntGrunts(outputPath string) (map[string]terragrunt.Config, error) {
	result := make(map[string]terragrunt.Config)

	for _, block := range s.Blocks {
		tgConfig, err := block.GenTerragruntGrunt()
		if err != nil {
			return nil, err
		}
		basePath := "."
		if s.Name == "" {
			basePath = fmt.Sprintf("%s/%s", basePath, s.Name)
		}
		tgConfigName := fmt.Sprintf("%s/%s", basePath, block.Name)
		if _, ok := result[tgConfigName]; ok {
			return nil, fmt.Errorf("duplicate Terragrunt configuration name: %s", tgConfigName)
		}
		result[tgConfigName] = tgConfig
	}

	for _, subSystems := range s.Systems {
		subResult, err := subSystems.GenTerragruntGrunts(outputPath)
		if err != nil {
			return nil, err
		}
		for k, v := range subResult {
			if _, ok := result[k]; ok {
				return nil, fmt.Errorf("duplicate Terragrunt configuration name: %s", k)
			}
			result[k] = v
		}
	}

	return result, nil
}
