package terragrunt

// DefaultTerragruntTemplate is the default Terragrunt configuration template used by Grunter.
const DefaultTerragruntTemplate = `{{ if .Dependencies }}# Dependencies{{end}}{{- range .Dependencies }}
dependency "{{.Name}}" {
  config_path  = "{{.ConfigPath}}"
  skip_outputs = {{.SkipOutputs}}
}
{{- end }}{{ if .Dependencies }}

{{end}}# Locals
locals {
  {{- range .LocalVariables }}
  {{.Name}} = {{.Value}}
  {{- end }}
}

# OpenTofu Configuration
terraform {
  {{- range .OpenTofu.BeforeHooks }}
  before_hook "{{.Name}}" {
    commands = [{{range .Commands}}"{{.}}",{{end}}]
    execute  = [{{range .Execute}}"{{.}}",{{end}}]
  }
  {{- end }}

  source = "{{.OpenTofu.Source}}"
}

# Include all settings from the root terragrunt.hcl file
include {
  path = find_in_parent_folders()
}

# Inputs to pass to the Terraform module
inputs = {
  {{- range $key, $value := .Inputs }}
  {{$key}} = {{$value}}
  {{- end }}
}
`
