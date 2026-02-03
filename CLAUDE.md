# CLAUDE.md

## All documents should be stored in the `docs` folder

## Don't write any comments in Go code unless specifically told to do so
## You could create TODO comments in Go code to remind yourself to update something later
## For existing comments in Go code keep them up to date in sync with the code

## Project structure should follow the standard Go project layout

## Build Commands (justfile)

- `just build` - Build the mdflux binary to `bin/mdflux`
- `just clean` - Remove build artifacts
- `just test` - Run tests
- `just coverage` - Run tests with coverage report
- `just deps` - Download and verify Go dependencies
- `just fmt` - Format Go code
- `just fetch-mermaid` - Download mermaid.min.js from CDN
- `just fetch-katex` - Build katex-embed tool and generate KaTeX CSS with embedded fonts
- `just fetch-assets` - Fetch all external assets (mermaid + katex)
- `just rebuild` - Full rebuild: clean, fetch all assets, and build