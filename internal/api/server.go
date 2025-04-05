// internal/api/server.go
package api

import (
	"net/http"

	"github.com/giovanni-gava/kube-validator/internal/parser"
	"github.com/giovanni-gava/kube-validator/internal/validator"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()

	r.POST("/validate", func(c *gin.Context) {
		yamlData, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		obj, err := parser.ParseYAMLToUnstructured(yamlData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := validator.CheckResources(obj)
		c.JSON(http.StatusOK, gin.H{"suggestions": result})
	})

	r.Run(":8080")
}
