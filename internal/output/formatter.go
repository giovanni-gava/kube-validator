// internal/output/formatter.go
package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/giovanni-gava/kube-validator/internal/validator"
)

func Print(results []validator.ValidationResult, format string) {
	switch format {
	case "json":
		printJSON(results)
	default:
		printPretty(results)
	}
}

func printJSON(results []validator.ValidationResult) {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println("Error generating JSON:", err)
		return
	}
	fmt.Println(string(data))
}

func printPretty(results []validator.ValidationResult) {
	for _, r := range results {
		msg := stripKnownPrefix(r.Message) // sempre limpa prefixos duplicados

		switch r.Level {
		case "error":
			color.New(color.FgRed).Printf("[ERROR] %s\n", msg)
		case "warning":
			color.New(color.FgYellow).Printf("[WARN]  %s\n", msg)
		case "suggestion":
			color.New(color.FgCyan).Printf("[SUGGESTION] %s\n", msg)
		case "info":
			color.New(color.FgGreen).Printf("[INFO] %s\n", msg)
		default:
			fmt.Println(msg)
		}

		if r.DocLink != "" {
			color.New(color.FgHiBlack).Printf("Doc: %s\n", r.DocLink)
		}
	}
}

func stripKnownPrefix(s string) string {
	// Também remove emojis para evitar duplicação como "💡 [SUGGESTION]"
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "💡")

	prefixes := []string{"[SUGGESTION]", "[ERROR]", "[INFO]", "[WARN]"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(s, prefix))
		}
	}
	return s
}
