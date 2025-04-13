package validator

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CheckContainersStructure(obj *unstructured.Unstructured) []ValidationResult {
	var results []ValidationResult
	kind := obj.GetKind()

	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if err != nil || !found {
		results = append(results, ValidationResult{
			Message: fmt.Sprintf("❌ [%s] spec.containers is missing or malformed – expected a list of containers under spec.containers[]", kind),
			Level:   "error",
		})
		return results
	}

	for _, c := range containers {
		container := c.(map[string]interface{})
		if _, ok := container["name"]; !ok {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("❌ [%s] container missing 'name' field in spec.containers[]", kind),
				Level:   "error",
			})
		}
		if _, ok := container["image"]; !ok {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("❌ [%s] container missing 'image' field in spec.containers[]", kind),
				Level:   "error",
			})
		}
	}

	return results
}
