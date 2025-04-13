package validator

import (
	"fmt"
	"strings"

	"github.com/giovanni-gava/kube-validator/internal/docs"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type FieldRule struct {
	Path       string
	Required   bool
	Type       string
	ExpectedAt string
	DocLink    string
}

func CheckFieldRules(obj *unstructured.Unstructured, rules []FieldRule) []ValidationResult {
	var results []ValidationResult
	kind := obj.GetKind()

	for _, rule := range rules {
		pathParts := strings.Split(rule.Path, ".")
		_, found, err := unstructured.NestedFieldNoCopy(obj.Object, pathParts...)

		docLink := rule.DocLink
		if docLink == "" {
			if link, ok := docs.FieldDocLinks[kind+"."+rule.Path]; ok {
				docLink = link
			}
		}

		if err != nil || !found {
			if rule.Required {
				results = append(results, ValidationResult{
					Message: fmt.Sprintf("❌ [%s] missing required field '%s'", kind, rule.Path),
					Level:   "error",
					DocLink: docLink,
				})
			}
			continue
		}

		if rule.ExpectedAt != "" && isFieldMisplaced(obj, rule.Path, rule.ExpectedAt) {
			results = append(results, ValidationResult{
				Message: fmt.Sprintf("❌ [%s] field '%s' was found at wrong location – expected under '%s'", kind, lastField(rule.Path), rule.ExpectedAt),
				Level:   "error",
				DocLink: docLink,
			})
			continue
		}
	}

	return results
}

func lastField(path string) string {
	parts := strings.Split(path, ".")
	return strings.TrimSuffix(parts[len(parts)-1], "[]")
}

func isFieldMisplaced(obj *unstructured.Unstructured, actualPath, expectedPath string) bool {
	containers, found, _ := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if !found {
		return false
	}
	for _, c := range containers {
		if container, ok := c.(map[string]interface{}); ok {
			if _, found := container["volumes"]; found {
				return true
			}
		}
	}
	return false
}

var PodFieldRules = []FieldRule{
	{
		Path:     "spec.containers",
		Required: true,
		Type:     "list",
	},
	{
		Path:       "spec.volumes",
		Required:   false,
		Type:       "list",
		ExpectedAt: "spec.volumes",
	},
}

func CheckPodStructure(obj *unstructured.Unstructured) []ValidationResult {
	if obj.GetKind() != "Pod" {
		return nil
	}
	return CheckFieldRules(obj, PodFieldRules)
}
