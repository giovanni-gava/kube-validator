package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	url := "https://raw.githubusercontent.com/kubernetes/kubernetes/master/api/openapi-spec/swagger.json"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var openapi map[string]interface{}
	if err := json.Unmarshal(body, &openapi); err != nil {
		panic(err)
	}

	definitions, ok := openapi["definitions"].(map[string]interface{})
	if !ok {
		panic("definitions not found in schema")
	}

	outputDir := filepath.Join("internal", "schema", "definitions")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(err)
	}

	for defName, defContent := range definitions {
		data, err := json.MarshalIndent(defContent, "", "  ")
		if err != nil {
			fmt.Printf("Failed to marshal %s: %v\n", defName, err)
			continue
		}

		// Sanitize filename: replace / and .
		filename := strings.ReplaceAll(defName, "/", "_")
		filename = strings.ReplaceAll(filename, ".", "_") + ".json"
		filePath := filepath.Join(outputDir, filename)

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			fmt.Printf("Failed to write %s: %v\n", filePath, err)
			continue
		}
	}

	fmt.Println("✅ All definitions extracted to:", outputDir)
}
