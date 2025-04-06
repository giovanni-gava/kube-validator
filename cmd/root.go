package cmd

import (
	"fmt"
	"os"

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

		parseResults := parser.ParseYAMLDocuments(data)
		var allValidations []validator.ValidationResult

		for _, res := range parseResults {
			if res.Error != nil {
				allValidations = append(allValidations, validator.ValidationResult{
					Message: fmt.Sprintf("❌ [SYNTAX] YAML syntax error: %v", res.Error),
					Level:   "error",
				})
				continue
			}

			obj := res.Object
			allValidations = append(allValidations, validator.CheckContainersStructure(obj)...)
			allValidations = append(allValidations, validator.CheckResources(obj)...)
			allValidations = append(allValidations, validator.CheckProbes(obj)...)
			allValidations = append(allValidations, validator.CheckImagePullPolicy(obj)...)
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
