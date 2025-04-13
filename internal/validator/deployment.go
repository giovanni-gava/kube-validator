package validator

import (
    "fmt"
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CheckDeploymentStructure(obj *unstructured.Unstructured) []ValidationResult {
    if obj.GetKind() != "Deployment" {
        return nil
    }

    var results []ValidationResult
    kind := obj.GetKind()

    name, found, _ := unstructured.NestedString(obj.Object, "metadata", "name")
    if !found || name == "" {
        results = append(results, ValidationResult{
            Message: fmt.Sprintf("❌ [%s] metadata.name is required", kind),
            Level:   "error",
        })
    }

    return results
}
