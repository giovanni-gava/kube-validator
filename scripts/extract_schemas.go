package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	inputPath := "internal/schema/deployment.json"
	outputRoot := "internal/schema"

	data, err := os.ReadFile(inputPath)
	if err != nil {
		panic(err)
	}

	var swagger map[string]interface{}
	err = json.Unmarshal(data, &swagger)
	if err != nil {
		panic(err)
	}

	definitions, ok := swagger["definitions"].(map[string]interface{})
	if !ok {
		panic("definitions not found in swagger.json")
	}

	for defName, def := range definitions {
		parts := strings.Split(defName, ".")
		if len(parts) < 5 {
			continue // skip malformed
		}
		group := parts[3]   // e.g., apps, core, networking
		version := parts[4] // e.g., v1
		kind := parts[len(parts)-1]

		// Special case: core group has no group path in Kubernetes (core/v1)
		if group == "core" {
			group = "core"
		}

		schema := map[string]interface{}{
			"$schema":     "http://json-schema.org/draft-07/schema#",
			"title":       defName,
			"type":        "object",
			"definitions": map[string]interface{}{defName: def},
			"$ref":        fmt.Sprintf("#/definitions/%s", defName),
		}

		outDir := filepath.Join(outputRoot, group, version)
		os.MkdirAll(outDir, 0755)

		filename := fmt.Sprintf("%s.json", kind)
		outputPath := filepath.Join(outDir, filename)

		bytes, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(outputPath, bytes, 0644)
		if err != nil {
			panic(err)
		}
		fmt.Printf("✓ Saved schema: %s\n", outputPath)
	}
}
