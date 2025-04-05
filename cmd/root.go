package cmd

import (
	"fmt"
	"os"

	"github.com/giovanni-gava/kube-validator/internal/parser"
	"github.com/giovanni-gava/kube-validator/internal/validator"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kube-validator",
	Short: "Valida arquivos YAML do Kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Você deve fornecer o caminho para o YAML")
			os.Exit(1)
		}

		data, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Println("Erro ao ler arquivo:", err)
			return
		}

		obj, err := parser.ParseYAMLToUnstructured(data)
		if err != nil {
			fmt.Println("Erro ao parsear:", err)
			return
		}

		result := validator.CheckResources(obj)
		for _, r := range result {
			fmt.Println("🔍", r)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
