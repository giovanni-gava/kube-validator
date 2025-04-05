package validator

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CheckProbes(obj *unstructured.Unstructured) []ValidationResult {
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
		name := container["name"]

		if _, ok, _ := unstructured.NestedMap(container, "livenessProbe"); !ok {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("⚠️ container %v is missing livenessProbe", name),
				Level:   "warning",
			})
		}

		if _, ok, _ := unstructured.NestedMap(container, "readinessProbe"); !ok {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("⚠️ container %v is missing readinessProbe", name),
				Level:   "warning",
			})
		}
	}

	return results
}
