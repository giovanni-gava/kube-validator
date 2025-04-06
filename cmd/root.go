package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/giovanni-gava/kube-validator/internal/parser"
	"github.com/giovanni-gava/kube-validator/internal/validator"
	"github.com/spf13/cobra"
)

var outputFormat string

var rootCmd = &cobra.Command{
	Use:   "kube-validator",
	Short: "Validate Kubernetes YAML files for correctness and best practices",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must provide the path to the YAML file")
			os.Exit(1)
		}

		data, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		// Step 1: Structural linting before parsing into K8s objects
		structuralResults, err := validator.LintRawYAMLFromBytes(data)
		var allValidations []validator.ValidationResult
		if err != nil {
			fmt.Println("Error linting YAML structure:", err)
		} else {
			allValidations = append(allValidations, structuralResults...)
		}

		// Step 2: Parse and validate Kubernetes semantics
		parseResults := parser.ParseYAMLDocuments(data)
		for _, res := range parseResults {
			if res.Error != nil {
				errMsg := res.Error.Error()
				hint := "General YAML formatting issue."
				docLink := "https://kubernetes.io/docs/concepts/configuration/overview/#understanding-yaml"

				if strings.Contains(errMsg, "did not find expected key") {
					hint = "Missing ':' or incorrect indentation. Example: ensure 'containers' is properly aligned under 'spec'."
				} else if strings.Contains(errMsg, "mapping values are not allowed here") {
					hint = "Check for misplaced colons or unquoted strings with special characters."
				} else if strings.Contains(errMsg, "cannot unmarshal") {
					hint = "Type mismatch or invalid structure. Ensure each field uses the correct format."
				}

				message := fmt.Sprintf("[ERROR] [SYNTAX] YAML syntax error: %v\nHint: %s", res.Error, hint)

				allValidations = append(allValidations, validator.ValidationResult{
					Message: message,
					Level:   "error",
					DocLink: docLink,
				})
				continue
			}

			obj := res.Object
			kind := obj.GetKind()

			switch kind {
			case "Deployment":
				allValidations = append(allValidations, validator.CheckDeploymentStructure(obj)...)
				allValidations = append(allValidations, validator.CheckResources(obj)...)
				allValidations = append(allValidations, validator.CheckProbes(obj)...)
				allValidations = append(allValidations, validator.CheckImagePullPolicy(obj)...)
			case "Pod":
				allValidations = append(allValidations, validator.CheckContainersStructure(obj)...)
				allValidations = append(allValidations, validator.CheckPodStructure(obj)...)
				allValidations = append(allValidations, validator.CheckResources(obj)...)
				allValidations = append(allValidations, validator.CheckProbes(obj)...)
				allValidations = append(allValidations, validator.CheckImagePullPolicy(obj)...)
			default:
				// Future support for other kinds
			}
		}

		// Only show success message if there are no errors or warnings
		hasErrorsOrWarnings := false
		for _, res := range allValidations {
			if res.Level == "error" || res.Level == "warning" {
				hasErrorsOrWarnings = true
				break
			}
		}

		if !hasErrorsOrWarnings {
			allValidations = append(allValidations, validator.ValidationResult{
				Message: "[INFO] Manifest is valid – no issues detected",
				Level:   "info",
			})
		}

		validator.PrintResults(allValidations, outputFormat)
	},
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
