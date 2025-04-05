package validator

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CheckResources(obj *unstructured.Unstructured) []ValidationResult {
	var results []ValidationResult

	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if err != nil || !found {
		return append(results, ValidationResult{
			Message: "❌ containers not found",
			Level:   "error",
		})
	}

	for _, c := range containers {
		container := c.(map[string]interface{})
		resources, found, _ := unstructured.NestedMap(container, "resources")
		if !found || len(resources) == 0 {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("⚠️ container %v is missing resources.requests/limits", container["name"]),
				Level:   "warning",
			})
		}
	}

	return results
}
