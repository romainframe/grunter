package block

import (
	"fmt"
	"strings"

	"github.com/romainframe/grunter/pkg/utils"
)

// K8sBeforeHooks defines pre-execution hooks for different Kubernetes environments.
// Currently, it supports a pre-execution hook for GCP.
var K8sBeforeHooks = map[string]BeforeHook{
	"gcp": {
		Name:     "gke-context",
		Commands: []string{"plan", "apply", "destroy"},
		Execute:  []string{"bash", "-c", "gcloud config set account ${local.email} && gcloud container clusters get-credentials ${local.cluster.locals.name} --region=${local.cluster.locals.region}  --project=${local.project.locals.slug}"},
	},
}

// K8sGruntBuilder initializes a GruntBuilder specifically for Kubernetes templates.
// It ensures the presence of required metadata and sets up necessary local configurations
// and pre-execution hooks based on the cloud environment.
var K8sGruntBuilder = GruntBuilder{
	// Matches returns true if the configuration template is intended for Kubernetes.
	Matches: func(c Block) bool {
		return strings.Contains(c.Template, "k8s")
	},
	// Build enriches the provided Block with Kubernetes-specific settings.
	// It validates the template, metadata, and cluster configuration, sets the cluster local,
	// and adds the appropriate pre-execution hook based on the cloud environment.
	Build: func(c Block) (Block, error) {
		if c.Template == "" {
			return c, ErrTemplateRequired
		}
		if len(c.Metadata) == 0 {
			return c, ErrMetadataRequired
		}

		// Validate and retrieve cluster metadata
		clusterMetadataValue, ok := c.Metadata["cluster"]
		if !ok || clusterMetadataValue == "" {
			return c, ErrMetadataKeyRequired("cluster")
		}

		// Find and set the cluster local configuration
		parentFolder := "services/k8s"
		clusterFile := "values.hcl"
		clusterFilePath, err := utils.FindFileInParentTarget(parentFolder, clusterMetadataValue, clusterFile, 50)
		if err != nil {
			return c, utils.WrapError(ErrFilePathNotFound(fmt.Sprintf("%s/%s/.../%s", parentFolder, clusterMetadataValue, clusterFile)), err)
		}
		c.Locals["cluster"] = fmt.Sprintf(`read_terragrunt_config("%s")`, clusterFilePath)

		// Determine the cloud environment type
		hcl, err := utils.GetHCLFromParent("cloud")
		if err != nil {
			return c, utils.WrapError(ErrInvalidFile("cloud.hcl"), err)
		}
		cloudType, ok := utils.Get(hcl, "locals.slug")
		if !ok || cloudType == "" {
			return c, utils.WrapError(ErrKeyRequired("cloud.locals.slug"), err)
		}

		// Add the appropriate pre-execution hook if it hasn't been added yet
		if !hasBeforeHook(c.BeforeHooks, fmt.Sprintf("%s-context", cloudType)) {
			k8sBeforeHook, ok := K8sBeforeHooks[cloudType]
			if !ok {
				return c, ErrBeforeHookNotFound(cloudType)
			}
			c.BeforeHooks = append(c.BeforeHooks, k8sBeforeHook)
		}

		return c, nil
	},
}

// hasBeforeHook checks if the specified hook name is already present in the given slice of BeforeHook.
func hasBeforeHook(hooks []BeforeHook, hookName string) bool {
	for _, hook := range hooks {
		if hook.Name == hookName {
			return true
		}
	}
	return false
}
