package validator

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CheckResources(obj *unstructured.Unstructured) []ValidationResult {
	var results []ValidationResult
	kind := obj.GetKind()

	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if err != nil || !found {
		return results // já será tratado no CheckContainersStructure
	}

	for _, c := range containers {
		container := c.(map[string]interface{})
		name := container["name"]

		resources, found, _ := unstructured.NestedMap(container, "resources")
		if !found || len(resources) == 0 {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("💡 [SUGGESTION] [%s] container '%v' is missing resource requests/limits – recommended for scheduling efficiency", kind, name),
				Level:   "suggestion",
			})
		}
	}

	return results
}
