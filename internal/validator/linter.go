package validator

import (
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"
)

type ASTFieldIssue struct {
	Path     string
	Message  string
	Level    string
	Expected string
}

func LintRawYAMLFromBytes(data []byte) ([]ValidationResult, error) {
	docs := strings.Split(string(data), "---")
	var all []ValidationResult

	for _, doc := range docs {
		if strings.TrimSpace(doc) == "" {
			continue
		}
		var m map[string]interface{}
		err := yaml.Unmarshal([]byte(doc), &m)
		if err != nil {
			all = append(all, ValidationResult{
				Message: fmt.Sprintf("[SYNTAX] YAML document parsing failed: %v", err),
				Level:   "error",
			})
			continue
		}
		results := LintRawYAMLStructure(m)
		all = append(all, results...)
	}

	return all, nil
}

func LintRawYAMLStructure(raw map[string]interface{}) []ValidationResult {
	var issues []ASTFieldIssue
	checkFieldsRecursive(raw, []string{}, &issues)

	var results []ValidationResult
	for _, issue := range issues {
		msg := fmt.Sprintf("[%s] %s", issue.Path, issue.Message)
		results = append(results, ValidationResult{
			Message: msg,
			Level:   issue.Level,
		})
	}
	return results
}

var knownCasing = map[string]string{
	"Image":           "image",
	"NAME":            "name",
	"STRATEGY":        "strategy",
	"TYPE":            "type",
	"liveness_probe":  "livenessProbe",
	"readiness_probe": "readinessProbe",
	"App":             "app",
}

func checkFieldsRecursive(obj interface{}, path []string, issues *[]ASTFieldIssue) {
	switch val := obj.(type) {
	case map[string]interface{}:
		for k, v := range val {
			fullPath := append(path, k)

			if expected, found := knownCasing[k]; found {
				*issues = append(*issues, ASTFieldIssue{
					Path:     strings.Join(fullPath, "."),
					Message:  fmt.Sprintf("Field '%s' is incorrectly cased. Expected '%s'", k, expected),
					Level:    "error",
					Expected: expected,
				})
			}

			if k == "volumes" && contains(path, "containers") {
				*issues = append(*issues, ASTFieldIssue{
					Path:    strings.Join(fullPath, "."),
					Message: "Field 'volumes' is misplaced inside 'containers'. Should be under spec.volumes",
					Level:   "error",
				})
			}

			checkFieldsRecursive(v, fullPath, issues)
		}
	case []interface{}:
		for idx, item := range val {
			checkFieldsRecursive(item, append(path, fmt.Sprintf("[%d]", idx)), issues)
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
