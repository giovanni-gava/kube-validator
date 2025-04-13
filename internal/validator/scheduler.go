package validator

import (
    "fmt"
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CheckAffinity(obj *unstructured.Unstructured) []ValidationResult {
    var results []ValidationResult
    kind := obj.GetKind()

    if _, found, _ := unstructured.NestedMap(obj.Object, "spec", "template", "spec", "affinity"); !found {
        results = append(results, ValidationResult{
            Message: fmt.Sprintf("💡 [SUGGESTION] [%s] consider defining affinity rules", kind),
            Level:   "suggestion",
        })
    }

    return results
}
