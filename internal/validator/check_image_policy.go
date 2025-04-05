package validator

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CheckImagePullPolicy(obj *unstructured.Unstructured) []string {
	var results []string
	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if err != nil || !found {
		return append(results, "❌ containers not found")
	}

	for _, c := range containers {
		container := c.(map[string]interface{})
		policy, found, _ := unstructured.NestedString(container, "imagePullPolicy")
		if !found || policy != "Always" {
			results = append(results, fmt.Sprintf("⚠️ container %v should use imagePullPolicy: Always", container["name"]))
		}
	}
	return results
}
