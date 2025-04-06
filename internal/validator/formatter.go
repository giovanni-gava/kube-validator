package validator

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
)

func PrintResults(results []ValidationResult, format string) {
	switch format {
	case "json":
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Println("Error while generating JSON:", err)
			return
		}
		fmt.Println(string(data))
	default:
		for _, r := range results {
			switch r.Level {
			case "error":
				color.New(color.FgRed).Printf("❌ [ERROR] %s\n", r.Message)
			case "warning":
				color.New(color.FgYellow).Printf("⚠️  [WARN]  %s\n", r.Message)
			case "suggestion":
				color.New(color.FgCyan).Printf("💡 [SUG]   %s\n", r.Message)
			default:
				fmt.Println(r.Message)
			}
		}
	}
}
