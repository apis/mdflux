package mermaid

import (
	_ "embed"
)

//go:embed assets/mermaid.min.js
var mermaidJS string

// MermaidJS returns the embedded mermaid.min.js source code.
func MermaidJS() string {
	return mermaidJS
}
