package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/giovanni-gava/kube-validator/internal/parser"
	"github.com/giovanni-gava/kube-validator/internal/validator"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func loadExample(t *testing.T, filename string) *unstructured.Unstructured {
	// Caminho: volta um nível para sair de /test e acessar /examples
	examplePath := filepath.Join("..", "examples", filename)

	data, err := os.ReadFile(examplePath)
	assert.NoError(t, err, "Erro ao ler o arquivo de exemplo")

	obj, err := parser.ParseYAMLToUnstructured(data)
	assert.NoError(t, err, "Erro ao fazer parse do YAML")

	return obj
}

func TestCheckResources(t *testing.T) {
	obj := loadExample(t, "pod-missing-resources.yaml")

	if obj == nil {
		t.Fatal("objeto de teste está nulo")
	}

	result := validator.CheckResources(obj)
	assert.NotEmpty(t, result)
	assert.Contains(t, result[0], "missing resources")
}

func TestCheckProbes(t *testing.T) {
	obj := loadExample(t, "pod-no-probes.yaml")

	if obj == nil {
		t.Fatal("objeto de teste está nulo")
	}

	result := validator.CheckProbes(obj)
	assert.NotEmpty(t, result)
	assert.Contains(t, result[0], "missing livenessProbe")
	assert.Contains(t, result[1], "missing readinessProbe")
}

func TestCheckImagePullPolicy(t *testing.T) {
	obj := loadExample(t, "pod-wrong-imagepolicy.yaml")

	if obj == nil {
		t.Fatal("objeto de teste está nulo")
	}

	result := validator.CheckImagePullPolicy(obj)
	assert.NotEmpty(t, result)
	assert.Contains(t, result[0], "should use imagePullPolicy: Always")
}
