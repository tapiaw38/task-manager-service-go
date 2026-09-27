package docs

import _ "embed"

//go:embed specs/openapi.yaml
var OpenAPISpec []byte

const OpenAPISpecContentType = "application/yaml"
