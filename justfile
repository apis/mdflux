# Set shell for Windows compatibility
set windows-shell := ["powershell.exe", "-NoLogo", "-Command"]

# Platform-specific executable suffix
exe_suffix := if os() == "windows" { ".exe" } else { "" }

# List available recipes
default:
    @just --list

# Build the mdflux binary
build:
    @echo "Building mdflux..."
    @go build -o bin/mdflux{{exe_suffix}} ./cmd/mdflux
    @echo "Build complete: bin/mdflux{{exe_suffix}}"

# Build and run with sample input
run: build
    @echo "Running mdflux..."
    @echo "# Sample\n\nThis is **bold** and *italic*." | ./bin/mdflux{{exe_suffix}} -i -

# Run with debug logging
run-debug: build
    @echo "Running mdflux with debug logging..."
    @echo "# Sample\n\nThis is **bold**." | ./bin/mdflux{{exe_suffix}} -i - -l debug

# Convert a file
convert input output="": build
    @echo "Converting {{input}}..."
    @./bin/mdflux{{exe_suffix}} -i {{input}} {{ if output != "" { "-o " + output } else { "" } }}

# Convert with full HTML document
convert-full input output="": build
    @echo "Converting {{input}} to full HTML document..."
    @./bin/mdflux{{exe_suffix}} -i {{input}} --html.full_document {{ if output != "" { "-o " + output } else { "" } }}

# Clean build artifacts
[unix]
clean:
    @echo "Cleaning..."
    @rm -rf bin/
    @go clean
    @echo "Clean complete"

[windows]
clean:
    @echo "Cleaning..."
    @if (Test-Path bin) { Remove-Item -Recurse -Force bin }
    @go clean
    @echo "Clean complete"

# Run tests
test:
    @echo "Running tests..."
    @go test -v ./...

# Run tests with coverage
coverage:
    @echo "Running tests with coverage..."
    @go test -coverprofile=coverage.out ./...
    @go tool cover -func=coverage.out
    @echo "Coverage report: coverage.out"

# Run tests with coverage and open HTML report
coverage-html: coverage
    @echo "Generating HTML coverage report..."
    @go tool cover -html=coverage.out -o coverage.html
    @echo "HTML report: coverage.html"

# Download and verify dependencies
deps:
    @echo "Downloading dependencies..."
    @go mod download
    @go mod verify
    @go mod tidy
    @echo "Dependencies ready"

# Format code
fmt:
    @echo "Formatting code..."
    @go fmt ./...
    @echo "Format complete"

# Run linter
[unix]
lint:
    @echo "Running linter..."
    @golangci-lint run || echo "Note: Install golangci-lint for linting"

[windows]
lint:
    @echo "Running linter..."
    @try { golangci-lint run } catch { echo "Note: Install golangci-lint for linting" }

# Show help
help: build
    @./bin/mdflux{{exe_suffix}} -?

# Install binary to GOPATH/bin
install:
    @echo "Installing mdflux..."
    @go install ./cmd/mdflux
    @echo "Installed to GOPATH/bin"

# Mermaid.js version to download
mermaid_version := "11.4.0"
mermaid_url := "https://cdn.jsdelivr.net/npm/mermaid@" + mermaid_version + "/dist/mermaid.min.js"
mermaid_dest := "web/assets/mermaid.min.js"

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

# KaTeX version to download
katex_version := "0.16.21"
katex_css_url := "https://cdn.jsdelivr.net/npm/katex@" + katex_version + "/dist/katex.min.css"
katex_font_base := "https://cdn.jsdelivr.net/npm/katex@" + katex_version + "/dist/fonts"
katex_dest := "web/assets/katex.min.css"

# Fetch katex.min.css from CDN and fix font URLs to use CDN
[unix]
fetch-katex:
    @echo "Fetching KaTeX CSS v{{katex_version}}..."
    @mkdir -p $(dirname {{katex_dest}})
    @curl -sL {{katex_css_url}} | sed 's|url(fonts/|url({{katex_font_base}}/|g' > {{katex_dest}}
    @echo "Downloaded to {{katex_dest}} ($(wc -c < {{katex_dest}} | tr -d ' ') bytes)"

[windows]
fetch-katex:
    @echo "Fetching KaTeX CSS v{{katex_version}}..."
    @New-Item -ItemType Directory -Force -Path (Split-Path {{katex_dest}}) | Out-Null
    @$css = (Invoke-WebRequest -Uri {{katex_css_url}}).Content; $css -replace 'url\(fonts/', 'url({{katex_font_base}}/' | Set-Content {{katex_dest}}
    @echo "Downloaded to {{katex_dest}}"

# Fetch all external assets
fetch-assets: fetch-mermaid fetch-katex
