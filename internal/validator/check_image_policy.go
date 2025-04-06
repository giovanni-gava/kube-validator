package validator

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CheckImagePullPolicy(obj *unstructured.Unstructured) []ValidationResult {
	var results []ValidationResult
	kind := obj.GetKind()

	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if err != nil || !found {
		return results // já tratado no CheckContainersStructure
	}

	for _, c := range containers {
		container := c.(map[string]interface{})
		name := container["name"]

		policy, found, _ := unstructured.NestedString(container, "imagePullPolicy")
		if !found || policy != "Always" {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("💡 [SUGGESTION] [%s] container '%v' should set imagePullPolicy: Always – ensures fresh images in CI workflows", kind, name),
				Level:   "suggestion",
			})
		}
	}

	return results
}
