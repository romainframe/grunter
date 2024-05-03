package terragrunt

const DefaultValuesTemplate = `locals { {{- range $key, $value := . }}
	{{$key}} = "{{$value}}"{{- end }}
}`
