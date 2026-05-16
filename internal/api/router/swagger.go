package router

import (
	"net/http"
	"todo/docs"
)

const swaggerUI = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Todo API Swagger</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
    <style>
      html, body {
        margin: 0;
        padding: 0;
        background: #f4f7fb;
      }
    </style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js" crossorigin></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: '/swagger/openapi.yaml',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [SwaggerUIBundle.presets.apis],
        layout: 'BaseLayout'
      });
    </script>
  </body>
</html>
`

func registerSwaggerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /swagger", serveSwaggerUI)
	mux.HandleFunc("GET /swagger/openapi.yaml", serveOpenAPI)
}

func serveSwaggerUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(swaggerUI))
}

func serveOpenAPI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(docs.OpenAPI)
}
