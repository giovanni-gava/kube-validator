package parser

import (
	"bytes"
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yamlserializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

type ParseResult struct {
	Object *unstructured.Unstructured
	Error  error
	Raw    string
}

var decUnstructured = yamlserializer.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

// ParseYAMLDocuments recebe YAMLs (com ou sem múltiplos documentos) e retorna ParseResults
func ParseYAMLDocuments(yamlContent []byte) []ParseResult {
	var results []ParseResult

	decoder := k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlContent), 4096)

	for {
		var raw map[string]interface{}
		err := decoder.Decode(&raw)
		if err == io.EOF {
			break
		}
		if err != nil {
			results = append(results, ParseResult{
				Object: nil,
				Error:  fmt.Errorf("erro de sintaxe: %v", err),
				Raw:    "",
			})
			continue
		}
		if len(raw) == 0 {
			continue
		}

		obj := &unstructured.Unstructured{Object: raw}
		results = append(results, ParseResult{
			Object: obj,
			Error:  nil,
			Raw:    "", // opcional
		})
	}

	return results
}
