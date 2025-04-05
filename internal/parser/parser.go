package parser

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

// Serializer para decodificar YAML para objetos do Kubernetes
var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

// ParseYAMLToUnstructured converte um arquivo YAML em um objeto genérico do Kubernetes (Unstructured)
func ParseYAMLToUnstructured(yamlContent []byte) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}
	_, _, err := decUnstructured.Decode(yamlContent, nil, obj)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear YAML: %v", err)
	}
	return obj, nil
}
