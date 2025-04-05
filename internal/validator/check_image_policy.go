package validator

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CheckImagePullPolicy(obj *unstructured.Unstructured) []ValidationResult {
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
		policy, found, _ := unstructured.NestedString(container, "imagePullPolicy")
		if !found || policy != "Always" {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("⚠️ container %v should use imagePullPolicy: Always", container["name"]),
				Level:   "suggestion",
			})
		}
	}

	return results
}
