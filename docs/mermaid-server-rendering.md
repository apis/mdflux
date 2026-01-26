# Mermaid Server-Side Rendering with chromedp

## Status: Implemented

This document describes the server-side mermaid rendering feature.

## Problem Statement (Solved)

The previous implementation used `go.abhg.dev/goldmark/mermaid` with `RenderModeClient`, which output raw mermaid code in `<div class="mermaid">` tags expecting browser-side JavaScript rendering. This caused:

1. **mermaid.js was NOT included** in the HTML templates
2. **chromedp doesn't wait for JS execution** before PDF capture
3. **Result: Mermaid diagrams didn't render in PDFs**

## Solution

Replaced client-side rendering with a custom goldmark extension that renders mermaid diagrams server-side using chromedp and an embedded mermaid.min.js file.

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Build/Setup Phase                        │
│  1. Download mermaid.min.js from CDN (just fetch-mermaid)   │
│  2. Place in internal/pkg/md-to-pdf/mermaid/assets/         │
│  3. Embed via //go:embed                                    │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    Runtime Flow                             │
│                                                             │
│  Markdown Input                                             │
│       ↓                                                     │
│  Goldmark parses mermaid code block                         │
│       ↓                                                     │
│  Custom mermaid.Extender (AST Transformer)                  │
│       ↓                                                     │
│  MermaidBlock node created                                  │
│       ↓                                                     │
│  Custom htmlRenderer.renderMermaidBlock()                   │
│       ↓                                                     │
│  ┌─────────────────────────────────────────┐                │
│  │  mermaid.Renderer (chromedp-based)      │                │
│  │  1. Launch headless Chrome (once)       │                │
│  │  2. Load HTML with embedded mermaid.js  │                │
│  │  3. Execute: mermaid.render(id, code)   │                │
│  │  4. Return SVG string                   │                │
│  └─────────────────────────────────────────┘                │
│       ↓                                                     │
│  SVG embedded inline in HTML output                         │
│       ↓                                                     │
│  PDF renders correctly (no JS needed)                       │
└─────────────────────────────────────────────────────────────┘
```

## Implementation

### Package Structure

```
internal/pkg/md-to-pdf/mermaid/
├── assets/
│   └── mermaid.min.js      # Downloaded from CDN via `just fetch-mermaid`
├── embed.go                 # //go:embed for mermaid.min.js
├── extension.go             # Custom goldmark extension (Extender, MermaidBlock, transformer, htmlRenderer)
└── renderer.go              # Chromedp-based Renderer
```

### Key Files

#### embed.go
Embeds mermaid.min.js into the binary using `//go:embed`.

#### extension.go
Custom goldmark extension with:
- `Extender` - implements `goldmark.Extender`
- `MermaidBlock` - custom AST node for mermaid diagrams
- `transformer` - converts mermaid fenced code blocks to MermaidBlock nodes
- `htmlRenderer` - renders MermaidBlock to SVG via the Renderer

#### renderer.go
Chromedp-based renderer that:
- Creates a temp HTML file with embedded mermaid.js
- Launches headless Chrome (reused across all diagrams)
- Executes `mermaid.render()` via JavaScript
- Returns the rendered SVG

### Usage

```bash
# Download mermaid.min.js (required before building)
just fetch-mermaid

# Build the binary
just build

# Convert markdown with mermaid diagrams
./bin/md-to-pdf -i document.md -o document.html
./bin/md-to-pdf -i document.md -o document.pdf -f pdf
```

### justfile Recipes

```just
# Mermaid.js version to download
mermaid_version := "11.4.0"
mermaid_url := "https://cdn.jsdelivr.net/npm/mermaid@" + mermaid_version + "/dist/mermaid.min.js"
mermaid_dest := "internal/pkg/md-to-pdf/mermaid/assets/mermaid.min.js"

# Fetch mermaid.min.js from CDN
[unix]
fetch-mermaid:
    @echo "Fetching mermaid.js v{{mermaid_version}}..."
    @mkdir -p $(dirname {{mermaid_dest}})
    @curl -sL {{mermaid_url}} -o {{mermaid_dest}}
    @echo "Downloaded to {{mermaid_dest}} ($(wc -c < {{mermaid_dest}} | tr -d ' ') bytes)"

[windows]
fetch-mermaid:
    @echo "Fetching mermaid.js v{{mermaid_version}}..."
    @New-Item -ItemType Directory -Force -Path (Split-Path {{mermaid_dest}}) | Out-Null
    @Invoke-WebRequest -Uri {{mermaid_url}} -OutFile {{mermaid_dest}}
    @echo "Downloaded to {{mermaid_dest}}"
```

## Benefits

| Benefit | Description |
|---------|-------------|
| **Single binary** | mermaid.js embedded, no external npm/node dependency |
| **PDF works correctly** | SVGs render without JavaScript execution |
| **HTML works correctly** | Inline SVGs work everywhere |
| **Reuses chromedp** | Already a dependency for PDF rendering |
| **Browser reuse** | Single Chrome instance for all diagrams in a document |
| **Error handling** | On render failure, falls back to code block with error comment |

## Supported Diagram Types

All mermaid diagram types are supported:
- Flowchart
- Sequence Diagram
- Class Diagram
- State Diagram
- Entity Relationship Diagram
- Gantt Chart
- Pie Chart
- Git Graph
- Mind Map
- Timeline
- Quadrant Chart
- User Journey
- And more...

## Future Improvements

- **Caching**: Cache rendered SVGs by content hash to avoid re-rendering identical diagrams
- **Parallel rendering**: Render multiple diagrams concurrently using goroutines
- **Theme support**: Pass mermaid theme configuration based on document theme (light/dark)
