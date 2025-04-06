package validator

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CheckProbes(obj *unstructured.Unstructured) []ValidationResult {
	var results []ValidationResult
	kind := obj.GetKind()

	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if err != nil || !found {
		return results // já tratado em CheckContainersStructure
	}

	for _, c := range containers {
		container := c.(map[string]interface{})
		name := container["name"]

		if _, ok, _ := unstructured.NestedMap(container, "livenessProbe"); !ok {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("💡 [SUGGESTION] [%s] container '%v' is missing livenessProbe – recommended for better availability", kind, name),
				Level:   "suggestion",
			})
		}
		if _, ok, _ := unstructured.NestedMap(container, "readinessProbe"); !ok {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("💡 [SUGGESTION] [%s] container '%v' is missing readinessProbe – recommended for graceful rollout", kind, name),
				Level:   "suggestion",
			})
		}
	}

	return results
}
