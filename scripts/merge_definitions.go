package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func fixRefs(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, v2 := range val {
			if k == "$ref" {
				if refStr, ok := v2.(string); ok && strings.HasSuffix(refStr, ".json") {
					name := strings.TrimSuffix(filepath.Base(refStr), ".json")
					name = strings.ReplaceAll(name, "_", ".")
					val[k] = "#/definitions/" + name
				}
			} else {
				val[k] = fixRefs(v2)
			}
		}
		return val
	case []interface{}:
		for i, item := range val {
			val[i] = fixRefs(item)
		}
		return val
	default:
		return val
	}
}

func main() {
	definitions := make(map[string]interface{})
	defDir := "internal/schema/definitions"

	filepath.WalkDir(defDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		file, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Failed to read %s: %v\n", path, err)
			return nil
		}

		var def map[string]interface{}
		if err := json.Unmarshal(file, &def); err != nil {
			fmt.Printf("Failed to parse %s: %v\n", path, err)
			return nil
		}

		// Reescreve os $ref
		def = fixRefs(def).(map[string]interface{})

		name := strings.TrimSuffix(filepath.Base(path), ".json")
		name = strings.ReplaceAll(name, "_", ".") // convert back to dot notation
		definitions[name] = def
		return nil
	})

	out := map[string]interface{}{
		"definitions": definitions,
	}

	outBytes, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		panic(err)
	}

	os.MkdirAll("internal/schema/merged", 0755)
	os.WriteFile("internal/schema/merged/openapi.json", outBytes, 0644)
	fmt.Println("✅ Merged schema written to internal/schema/merged/openapi.json")
}
