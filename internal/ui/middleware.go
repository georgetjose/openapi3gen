package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterSwaggerUI(r *gin.Engine, _ string) {
	r.GET("/swagger", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, swaggerHTML)
	})
}

// Basic embedded Swagger UI HTML pointing to the JSON file
const swaggerHTML = `
<!DOCTYPE html>
<html>
<head>
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script>
    const ui = SwaggerUIBundle({
    url: '/swagger/openapi.json',
    dom_id: '#swagger-ui',
    });
  </script>
</body>
</html>
`
