package docs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/task-manager-service-go/docs"
)

const swaggerUI = `<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Task Manager Service</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.SwaggerUIBundle({
        url: './docs/openapi.yaml',
        dom_id: '#swagger-ui',
        docExpansion: 'list',
        defaultModelsExpandDepth: 1,
      });
    };
  </script>
</body>
</html>`

func NewUIHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUI))
	}
}

func NewSpecHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, docs.OpenAPISpecContentType, docs.OpenAPISpec)
	}
}
