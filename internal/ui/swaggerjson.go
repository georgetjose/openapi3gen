package ui

import (
	"net/http"

	"github.com/georgetjose/openapi3gen/internal/generator"
	"github.com/gin-gonic/gin"
)

// RegisterSwaggerJSONHandler mounts GET /swagger/openapi.json
func RegisterSwaggerJSONHandler(r *gin.Engine, openapi *generator.OpenAPI) {
	r.GET("/swagger/openapi.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, openapi)
	})
}
