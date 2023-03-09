package docs

import (
	_ "embed"
)

//go:embed swagger/swagger.json
var SwaggerJSON string
