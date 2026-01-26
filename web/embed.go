package web

import (
	"embed"
)

//go:embed templates
var TemplateFS embed.FS

//go:embed assets/mermaid.min.js
var MermaidJS string
