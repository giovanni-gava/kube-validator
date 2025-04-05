package validator

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CheckResources(obj *unstructured.Unstructured) []string {
	var results []string
	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if err != nil || !found {
		return append(results, "❌ containers not found")
	}

	for _, c := range containers {
		container := c.(map[string]interface{})
		resources, found, _ := unstructured.NestedMap(container, "resources")
		if !found || len(resources) == 0 {
			results = append(results, fmt.Sprintf("⚠️ container %v is missing resources.requests/limits", container["name"]))
		}
	}

	return results
}
