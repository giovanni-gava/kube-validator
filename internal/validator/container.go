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
			Message: fmt.Sprintf("❌ [%s] spec.containers is missing or malformed – expected a list of containers", kind),
			Level:   "error",
		})
		return results
	}

	for _, c := range containers {
		container := c.(map[string]interface{})
		nameVal, _ := container["name"]
		name := fmt.Sprintf("%v", nameVal)

		if _, ok := container["image"]; !ok {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("❌ [%s] container '%s' missing 'image' field", kind, name),
				Level:   "error",
			})
		}
	}

	return results
}

func CheckProbes(obj *unstructured.Unstructured) []ValidationResult {
	var results []ValidationResult
	kind := obj.GetKind()

	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if err != nil || !found {
		return results
	}

	for _, c := range containers {
		container := c.(map[string]interface{})
		nameVal, _ := container["name"]
		name := fmt.Sprintf("%v", nameVal)

		if _, ok, _ := unstructured.NestedMap(container, "livenessProbe"); !ok {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("💡 [SUGGESTION] [%s] container '%s' is missing livenessProbe", kind, name),
				Level:   "suggestion",
			})
		}
		if _, ok, _ := unstructured.NestedMap(container, "readinessProbe"); !ok {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("💡 [SUGGESTION] [%s] container '%s' is missing readinessProbe", kind, name),
				Level:   "suggestion",
			})
		}
	}

	return results
}

func CheckImagePullPolicy(obj *unstructured.Unstructured) []ValidationResult {
	var results []ValidationResult
	kind := obj.GetKind()

	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if err != nil || !found {
		return results
	}

	for _, c := range containers {
		container := c.(map[string]interface{})
		nameVal, _ := container["name"]
		name := fmt.Sprintf("%v", nameVal)

		policy, found, _ := unstructured.NestedString(container, "imagePullPolicy")
		if !found || policy != "Always" {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("💡 [SUGGESTION] [%s] container '%s' should set imagePullPolicy: Always", kind, name),
				Level:   "suggestion",
			})
		}
	}

	return results
}

func CheckResources(obj *unstructured.Unstructured) []ValidationResult {
	var results []ValidationResult
	kind := obj.GetKind()

	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if err != nil || !found {
		return results
	}

	for _, c := range containers {
		container := c.(map[string]interface{})
		nameVal, _ := container["name"]
		name := fmt.Sprintf("%v", nameVal)

		resources, found, _ := unstructured.NestedMap(container, "resources")
		if !found || len(resources) == 0 {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("💡 [SUGGESTION] [%s] container '%s' is missing resource requests/limits", kind, name),
				Level:   "suggestion",
			})
		}
	}

	return results
}
