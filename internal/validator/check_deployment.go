package validator

import (
	"fmt"

	"github.com/giovanni-gava/kube-validator/internal/docs"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CheckDeploymentStructure(obj *unstructured.Unstructured) []ValidationResult {
	if obj.GetKind() != "Deployment" {
		return nil
	}

	var results []ValidationResult
	kind := obj.GetKind()

	// metadata.name
	name, found, _ := unstructured.NestedString(obj.Object, "metadata", "name")
	if !found || name == "" {
		results = append(results, ValidationResult{
			Message: fmt.Sprintf("❌ [%s] metadata.name is required", kind),
			Level:   "error",
			DocLink: "https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names",
		})
	}

	// selector.matchLabels
	if _, found, _ := unstructured.NestedMap(obj.Object, "spec", "selector", "matchLabels"); !found {
		results = append(results, ValidationResult{
			Message: fmt.Sprintf("❌ [%s] spec.selector.matchLabels is required", kind),
			Level:   "error",
			DocLink: docs.FieldDocLinks["Deployment.spec.selector"],
		})
	}

	// selector vs labels match
	selectorLabels, _, _ := unstructured.NestedMap(obj.Object, "spec", "selector", "matchLabels")
	templateLabels, _, _ := unstructured.NestedMap(obj.Object, "spec", "template", "metadata", "labels")
	if !mapsMatch(selectorLabels, templateLabels) {
		results = append(results, ValidationResult{
			Message: fmt.Sprintf("❌ [%s] spec.template.metadata.labels must match spec.selector.matchLabels", kind),
			Level:   "error",
			DocLink: docs.FieldDocLinks["Deployment.spec.selector"],
		})
	}

	// replicas >= 1
	replicas, found, _ := unstructured.NestedInt64(obj.Object, "spec", "replicas")
	if found && replicas < 1 {
		results = append(results, ValidationResult{
			Message: fmt.Sprintf("💡 [SUGGESTION] [%s] replicas is less than 1, consider setting at least 1 for availability", kind),
			Level:   "suggestion",
			DocLink: docs.FieldDocLinks["Deployment.spec.replicas"],
		})
	}

	// rollingUpdate strategy config
	strategyType, _, _ := unstructured.NestedString(obj.Object, "spec", "strategy", "type")
	if strategyType == "RollingUpdate" {
		_, surgeFound, _ := unstructured.NestedFieldNoCopy(obj.Object, "spec", "strategy", "rollingUpdate", "maxSurge")
		_, unavailableFound, _ := unstructured.NestedFieldNoCopy(obj.Object, "spec", "strategy", "rollingUpdate", "maxUnavailable")
		if !surgeFound || !unavailableFound {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("💡 [SUGGESTION] [%s] rollingUpdate strategy should define maxSurge and maxUnavailable", kind),
				Level:   "suggestion",
				DocLink: docs.FieldDocLinks["Deployment.spec.strategy.rollingUpdate.maxUnavailable"],
			})
		}
	}

	// affinity
	if _, found, _ := unstructured.NestedMap(obj.Object, "spec", "template", "spec", "affinity"); !found {
		results = append(results, ValidationResult{
			Message: fmt.Sprintf("💡 [SUGGESTION] [%s] consider defining affinity rules for pod placement", kind),
			Level:   "suggestion",
			DocLink: docs.FieldDocLinks["Deployment.spec.template.spec.affinity"],
		})
	}

	// tolerations
	if _, found, _ := unstructured.NestedSlice(obj.Object, "spec", "template", "spec", "tolerations"); !found {
		results = append(results, ValidationResult{
			Message: fmt.Sprintf("💡 [SUGGESTION] [%s] consider defining tolerations for scheduling flexibility", kind),
			Level:   "suggestion",
			DocLink: docs.FieldDocLinks["Deployment.spec.template.spec.tolerations"],
		})
	}

	// nodeSelector
	if _, found, _ := unstructured.NestedMap(obj.Object, "spec", "template", "spec", "nodeSelector"); !found {
		results = append(results, ValidationResult{
			Message: fmt.Sprintf("💡 [SUGGESTION] [%s] consider defining nodeSelector for targeted scheduling", kind),
			Level:   "suggestion",
			DocLink: docs.FieldDocLinks["Deployment.spec.template.spec.nodeSelector"],
		})
	}

	// securityContext.runAsNonRoot
	containers, found, _ := unstructured.NestedSlice(obj.Object, "spec", "template", "spec", "containers")
	if found {
		for _, c := range containers {
			container := c.(map[string]interface{})
			name := container["name"].(string)
			securityContext, found := container["securityContext"].(map[string]interface{})
			if !found || securityContext["runAsNonRoot"] != true {
				results = append(results, ValidationResult{
					Message: fmt.Sprintf("💡 [SUGGESTION] [Container %s] should have runAsNonRoot: true", name),
					Level:   "suggestion",
					DocLink: docs.FieldDocLinks["Deployment.spec.template.spec.containers[].securityContext.runAsNonRoot"],
				})
			}
		}
	}

	// annotations.kubectl.kubernetes.io/restartedAt
	annotations, found, _ := unstructured.NestedStringMap(obj.Object, "spec", "template", "metadata", "annotations")
	if found {
		if _, ok := annotations["kubectl.kubernetes.io/restartedAt"]; !ok {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("💡 [SUGGESTION] [%s] add annotation kubectl.kubernetes.io/restartedAt for manual restart tracking", kind),
				Level:   "suggestion",
				DocLink: docs.FieldDocLinks["Deployment.spec.template.metadata.annotations"],
			})
		}
	}

	return results
}

func mapsMatch(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
