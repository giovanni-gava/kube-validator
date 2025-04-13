// cmd/root.go
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/giovanni-gava/kube-validator/internal/output"
	"github.com/giovanni-gava/kube-validator/internal/parser"
	"github.com/giovanni-gava/kube-validator/internal/validator"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var outputFormat string

var rootCmd = &cobra.Command{
	Use:   "kube-validator",
	Short: "Validate Kubernetes YAML files for correctness and best practices",
	Run:   runValidator,
}

func init() {
	rootCmd.Flags().StringVarP(&outputFormat, "output", "o", "pretty", "Output format: pretty | json")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runValidator(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("You must provide the path to the YAML file")
		os.Exit(1)
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Step 1: Structural linting
	structuralResults, err := validator.LintRawYAMLFromBytes(data)
	var allValidations []validator.ValidationResult
	if err != nil {
		fmt.Println("Error linting YAML structure:", err)
	} else {
		allValidations = append(allValidations, structuralResults...)
	}

	// Step 2: Semantic validation
	parseResults := parser.ParseYAMLDocuments(data)
	for _, res := range parseResults {
		if res.Error != nil {
			appendSyntaxError(&allValidations, res.Error)
			continue
		}

		obj := res.Object
		kind := obj.GetKind()
		apiVersion := obj.GetAPIVersion()

		// Schema validation
		schemaResults := validator.ValidateWithSchema(obj.Object, kind, apiVersion)
		allValidations = append(allValidations, schemaResults...)

		// Custom semantic validators
		allValidations = append(allValidations, applyValidatorsForKind(kind, obj)...)
	}

	if !hasIssues(allValidations) {
		allValidations = append(allValidations, validator.ValidationResult{
			Message: "[INFO] Manifest is valid – no issues detected",
			Level:   "info",
		})
	}

	output.Print(allValidations, outputFormat)
}

func applyValidatorsForKind(kind string, obj *unstructured.Unstructured) []validator.ValidationResult {
	switch kind {
	case "Deployment":
		return applyValidators(obj,
			validator.CheckDeploymentStructure,
			validator.CheckContainersStructure,
			validator.CheckProbes,
			validator.CheckImagePullPolicy,
			validator.CheckResources,
			validator.CheckAffinity,
		)
	case "Pod":
		return applyValidators(obj,
			validator.CheckContainersStructure,
			validator.CheckProbes,
			validator.CheckImagePullPolicy,
			validator.CheckResources,
		)
	default:
		return nil
	}
}

func applyValidators(obj *unstructured.Unstructured, validators ...func(*unstructured.Unstructured) []validator.ValidationResult) []validator.ValidationResult {
	var results []validator.ValidationResult
	for _, v := range validators {
		results = append(results, v(obj)...)
	}
	return results
}

func appendSyntaxError(results *[]validator.ValidationResult, err error) {
	hint := "General YAML formatting issue."
	docLink := "https://kubernetes.io/docs/concepts/configuration/overview/#understanding-yaml"
	errMsg := err.Error()

	if strings.Contains(errMsg, "did not find expected key") {
		hint = "Missing ':' or incorrect indentation."
	} else if strings.Contains(errMsg, "mapping values are not allowed here") {
		hint = "Check for misplaced colons or unquoted strings."
	} else if strings.Contains(errMsg, "cannot unmarshal") {
		hint = "Type mismatch or invalid structure."
	}

	message := fmt.Sprintf("[ERROR] [SYNTAX] YAML syntax error: %v\nHint: %s", errMsg, hint)
	*results = append(*results, validator.ValidationResult{
		Message: message,
		Level:   "error",
		DocLink: docLink,
	})
}

func hasIssues(results []validator.ValidationResult) bool {
	for _, res := range results {
		if res.Level == "error" || res.Level == "warning" {
			return true
		}
	}
	return false
}
