package validator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func ValidateWithSchema(obj map[string]interface{}, kind string, apiVersion string) []ValidationResult {
	var results []ValidationResult

	schemaPath := "internal/schema/merged/openapi.json"
	schemaFile, err := os.ReadFile(schemaPath)
	if err != nil {
		return []ValidationResult{{
			Level:   "error",
			Message: fmt.Sprintf("Failed to read merged schema: %v", err),
		}}
	}

	var schemaData map[string]interface{}
	if err := json.Unmarshal(schemaFile, &schemaData); err != nil {
		return []ValidationResult{{
			Level:   "error",
			Message: fmt.Sprintf("Failed to parse schema: %v", err),
		}}
	}

	definitionKey := fmt.Sprintf("io.k8s.api.%s.%s", apiGroupVersion(apiVersion), kind)
	definition, found := schemaData["definitions"].(map[string]interface{})[definitionKey]
	if !found {
		return []ValidationResult{{
			Level:   "error",
			Message: fmt.Sprintf("No schema definition found for %s (%s)", kind, apiVersion),
		}}
	}

	required, ok := definition.(map[string]interface{})["required"].([]interface{})
	if !ok || required == nil {
		required = []interface{}{} // Corrige o erro: Garante sempre um array
	}

	specificSchema := map[string]interface{}{
		"$schema":     "http://json-schema.org/draft-07/schema#",
		"title":       definitionKey,
		"type":        "object",
		"properties":  definition.(map[string]interface{})["properties"],
		"required":    required, // Usa a correção aqui
		"definitions": schemaData["definitions"],
	}

	specificSchemaJSON, err := json.Marshal(specificSchema)
	if err != nil {
		return []ValidationResult{{
			Level:   "error",
			Message: fmt.Sprintf("Failed to encode specific schema: %v", err),
		}}
	}

	schemaLoader := gojsonschema.NewSchemaLoader()
	schemaLoader.Validate = true

	schemaDoc := gojsonschema.NewBytesLoader(specificSchemaJSON)
	schema, err := schemaLoader.Compile(schemaDoc)
	if err != nil {
		return []ValidationResult{{
			Level:   "error",
			Message: fmt.Sprintf("Failed to compile schema: %v", err),
		}}
	}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return []ValidationResult{{
			Level:   "error",
			Message: fmt.Sprintf("Failed to encode manifest to JSON: %v", err),
		}}
	}

	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)
	res, err := schema.Validate(documentLoader)
	if err != nil {
		return []ValidationResult{{
			Level:   "error",
			Message: fmt.Sprintf("Schema validation error: %v", err),
		}}
	}

	if !res.Valid() {
		for _, e := range res.Errors() {
			results = append(results, ValidationResult{
				Level:   "error",
				Message: fmt.Sprintf("[%s] %s", e.Field(), e.Description()),
			})
		}
	}

	return results
}

// Função já corrigida anteriormente
func apiGroupVersion(apiVersion string) string {
	if apiVersion == "v1" {
		return "core.v1"
	}
	return strings.ReplaceAll(apiVersion, "/", ".")
}
