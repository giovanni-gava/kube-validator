package validator

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CheckProbes(obj *unstructured.Unstructured) []string {
	var results []string
	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if err != nil || !found {
		return append(results, "❌ containers not found")
	}

	for _, c := range containers {
		container := c.(map[string]interface{})
		name := container["name"]

		if _, ok, _ := unstructured.NestedMap(container, "livenessProbe"); !ok {
			results = append(results, fmt.Sprintf("⚠️ container %v is missing livenessProbe", name))
		}
		if _, ok, _ := unstructured.NestedMap(container, "readinessProbe"); !ok {
			results = append(results, fmt.Sprintf("⚠️ container %v is missing readinessProbe", name))
		}
	}
	return results
}
