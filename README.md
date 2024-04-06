# Grunter

`'grʌntər`

> A worker, often in a technical, creative, or professional field, who performs the repetitive, tedious, or less glamorous tasks necessary for a project or operation. This may involve any work deemed necessary but not necessarily requiring high-level decision-making or creativity.

[![Go Report Card](https://goreportcard.com/badge/github.com/romainframe/grunter)](https://goreportcard.com/report/github.com/romainframe/grunter)
[![Go Reference](https://pkg.go.dev/badge/github.com/romainframe/grunter.svg)](https://pkg.go.dev/github.com/romainframe/grunter)
[![GitHub](https://img.shields.io/github/license/romainframe/grunter)]

Grunter is a simple Terragrunt file generator. It generates a `terragrunt.hcl` file from your `config.yaml` file in the current directory.

It enables you to define your Terragrunt configuration in a simple YAML file, which is easier to read and maintain than the HCL format. It saves you from writing the boilerplate code for each module and keeps your IaaC configuration even DRYer.

It is meant to be used in conjunction with Terragrunt, a thin wrapper for Terraform that provides extra tools for keeping your configurations DRY and managing remote state.

The configuration format is defined [here](./pkg/config/config.go).

## Installation

```bash
go install github.com/romainframe/grunter@latest
```

## Usage

```bash
grunter gen
```

## Example

Given the following `config.yaml` file:

```yaml
template: modules/k8s/namespace/default
metadata:
  cluster: main
dependencies:
  - name: cluster
    path: ../../../cluster
    pathType: relative
    withOutputs: true
  - name: env
    path: "${local.env}/example-${local.stack_namespace}"
    pathType: root
inputs:
  name: local.values.locals.name
  namespace: local.values.locals.namespace
```

The following `terragrunt.hcl` file will be generated:

```hcl
# Dependencies
dependency "cluster" {
  config_path  = "../../../cluster"
  skip_outputs = false
}
dependency "env" {
  config_path  = "${local.env}/example-${local.stack_namespace}"
  skip_outputs = true
}

# Locals
locals {
  cluster = read_terragrunt_config("../../../services/k8s/main/values.hcl")
  # grunted locals = begin
  email         = get_env("TF_VAR_EMAIL", "")
  project       = read_terragrunt_config(find_in_parent_folders("project.hcl"))
  template_root = get_env("TF_VAR_TEMPLATE_ROOT", "")
  values        = read_terragrunt_config("values.hcl")
  # grunted locals = end
}

# OpenTofu Configuration
terraform {
  before_hook "gke-context" {
    commands = ["plan", "apply", "destroy", ]
    execute  = ["bash", "-c", "gcloud config set account ${local.email} && gcloud container clusters get-credentials ${local.cluster.locals.name} --region=${local.cluster.locals.region}  --project=${local.project.locals.slug}", ]
  }

  source = "${local.template_root}//modules/k8s/namespace/default"
}

# Include all settings from the root terragrunt.hcl file
include {
  path = find_in_parent_folders()
}

# Inputs to pass to the Terraform module
inputs = {
  name      = local.values.locals.name
  namespace = local.values.locals.namespace
}
```

## Contributors

Without the contributions from these fine folks, this project would be a total dud!

<a href="https://github.com/romainframe/grunter/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=romainframe/grunter" />
</a>

---

© 2024 Romain Untereiner. All materials licensed under [Apache v2.0](http://www.apache.org/licenses/LICENSE-2.0)
