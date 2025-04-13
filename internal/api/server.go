package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giovanni-gava/kube-validator/internal/output"
	"github.com/giovanni-gava/kube-validator/internal/parser"
	"github.com/giovanni-gava/kube-validator/internal/validator"
)

func StartServer() {
	r := gin.Default()

	r.POST("/validate", func(c *gin.Context) {
		yamlData, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		// Parse múltiplos documentos YAML
		parseResults := parser.ParseYAMLDocuments(yamlData)

		var allValidations []validator.ValidationResult
		for _, res := range parseResults {
			if res.Error != nil {
				allValidations = append(allValidations, validator.ValidationResult{
					Message: "[ERROR] YAML syntax error: " + res.Error.Error(),
					Level:   "error",
				})
				continue
			}

			obj := res.Object
			kind := obj.GetKind()
			apiVersion := obj.GetAPIVersion()

			allValidations = append(allValidations, validator.ValidateWithSchema(obj.Object, kind, apiVersion)...)

			switch kind {
			case "Deployment":
				allValidations = append(allValidations, validator.CheckDeploymentStructure(obj)...)
				allValidations = append(allValidations, validator.CheckResources(obj)...)
				allValidations = append(allValidations, validator.CheckProbes(obj)...)
				allValidations = append(allValidations, validator.CheckImagePullPolicy(obj)...)
			case "Pod":
				allValidations = append(allValidations, validator.CheckContainersStructure(obj)...)
				allValidations = append(allValidations, validator.CheckPodStructure(obj)...)
				allValidations = append(allValidations, validator.CheckResources(obj)...)
				allValidations = append(allValidations, validator.CheckProbes(obj)...)
				allValidations = append(allValidations, validator.CheckImagePullPolicy(obj)...)
			}
		}

		// Suporte para output=pretty|json via query param
		format := c.DefaultQuery("output", "json")

		if format == "pretty" {
			output.Print(allValidations, "pretty")
			c.JSON(http.StatusOK, gin.H{"status": "printed to stdout"})
		} else {
			c.JSON(http.StatusOK, allValidations)
		}
	})

	r.Run(":8080")
}
