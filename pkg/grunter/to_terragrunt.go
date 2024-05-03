package grunter

import "github.com/romainframe/grunter/pkg/terragrunt"

func (o Object) GenTerragruntGrunts(outputPath string) (map[string]terragrunt.Config, error) {

	switch o.Kind {
	case ObjectKindBlock:
		tfConfig, err := o.block.GenTerragruntGrunt()
		if err != nil {
			return nil, err
		}
		return map[string]terragrunt.Config{outputPath: tfConfig}, nil
	case ObjectKindSystem:
		return o.system.GenTerragruntGrunts(outputPath)
	}
	return nil, nil
}
