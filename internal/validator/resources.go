package validator

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CheckResources(obj *unstructured.Unstructured) []string {
	var suggestions []string

	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if err != nil || !found {
		return append(suggestions, "containers not found")
	}

	for _, c := range containers {
		containerMap := c.(map[string]interface{})

		resources, found, _ := unstructured.NestedMap(containerMap, "resources")
		if !found || resources == nil {
			suggestions = append(suggestions, fmt.Sprintf("container %v is missing resources.requests/limits", containerMap["name"]))
		}
	}

	return suggestions
}
