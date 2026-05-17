// Package docs provide API-documentation
package docs

import _ "embed"

// OpenAPI содержит сгенерированную спецификацию API.
//
//go:embed openapi.yaml
var OpenAPI []byte
