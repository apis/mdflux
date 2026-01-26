package mermaid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"md-to-pdf/web"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
)

// Renderer implements the mermaid.ServerRenderer interface using chromedp.
// It renders mermaid diagram code to SVG using a headless Chrome browser.
type Renderer struct {
	ctx         context.Context
	cancel      context.CancelFunc
	allocCancel context.CancelFunc
	mu          sync.Mutex
	initialized bool
	chromePath  string
	tempDir     string
	htmlPath    string
}

// NewRenderer creates a new mermaid renderer.
// The chromePath parameter is optional; if empty, chromedp will auto-detect Chrome.
func NewRenderer(chromePath string) *Renderer {
	return &Renderer{
		chromePath: chromePath,
	}
}

// init lazily initializes the browser instance on first render.
func (r *Renderer) init() error {
	if r.initialized {
		return nil
	}

	// Create temp directory for HTML files
	tempDir, err := os.MkdirTemp("", "mermaid-renderer-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	r.tempDir = tempDir

	htmlTmplContent, err := web.TemplateFS.ReadFile("templates/mermaid.html")
	if err != nil {
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			log.Error().Err(removeErr).Str("path", tempDir).Msg("failed to remove temp directory")
		}
		return fmt.Errorf("failed to read mermaid HTML template: %w", err)
	}
	htmlTmpl, err := template.New("mermaid").Parse(string(htmlTmplContent))
	if err != nil {
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			log.Error().Err(removeErr).Str("path", tempDir).Msg("failed to remove temp directory")
		}
		return fmt.Errorf("failed to parse mermaid HTML template: %w", err)
	}
	var htmlBuf bytes.Buffer
	if err := htmlTmpl.Execute(&htmlBuf, map[string]string{"MermaidJS": web.MermaidJS}); err != nil {
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			log.Error().Err(removeErr).Str("path", tempDir).Msg("failed to remove temp directory")
		}
		return fmt.Errorf("failed to execute mermaid HTML template: %w", err)
	}
	htmlContent := htmlBuf.String()

	r.htmlPath = filepath.Join(tempDir, "mermaid.html")
	if err := os.WriteFile(r.htmlPath, []byte(htmlContent), 0644); err != nil {
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			log.Error().Err(removeErr).Str("path", tempDir).Msg("failed to remove temp directory")
		}
		return fmt.Errorf("failed to write mermaid HTML: %w", err)
	}

	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.NoSandbox,
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("no-zygote", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	if r.chromePath != "" {
		allocOpts = append(allocOpts, chromedp.ExecPath(r.chromePath))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	ctx, cancel := chromedp.NewContext(allocCtx)

	// Navigate to the mermaid HTML file to load mermaid.js
	fileURL := "file://" + r.htmlPath
	if err := chromedp.Run(ctx, chromedp.Navigate(fileURL), chromedp.WaitReady("body")); err != nil {
		allocCancel()
		cancel()
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			log.Error().Err(removeErr).Str("path", tempDir).Msg("failed to remove temp directory")
		}
		return fmt.Errorf("failed to initialize browser: %w", err)
	}

	r.ctx = ctx
	r.cancel = cancel
	r.allocCancel = allocCancel
	r.initialized = true

	return nil
}

// Render renders the given mermaid code to SVG.
// This method implements the mermaid.ServerRenderer interface.
func (r *Renderer) Render(code string) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.init(); err != nil {
		return nil, err
	}

	// Escape the code for safe embedding in JavaScript
	codeJSON, err := json.Marshal(code)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal mermaid code: %w", err)
	}

	jsTmplContent, err := web.TemplateFS.ReadFile("templates/mermaid-render.js")
	if err != nil {
		return nil, fmt.Errorf("failed to read mermaid render JS template: %w", err)
	}
	jsTmpl, err := template.New("render").Parse(string(jsTmplContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse mermaid render JS template: %w", err)
	}
	var jsBuf bytes.Buffer
	if err := jsTmpl.Execute(&jsBuf, map[string]string{"Code": string(codeJSON)}); err != nil {
		return nil, fmt.Errorf("failed to execute mermaid render JS template: %w", err)
	}
	renderScript := jsBuf.String()

	var result map[string]interface{}

	err = chromedp.Run(r.ctx,
		chromedp.Evaluate(renderScript, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("chromedp render failed: %w", err)
	}

	// Check for error
	if errVal, ok := result["error"]; ok && errVal != nil {
		if errStr, ok := errVal.(string); ok && errStr != "" {
			return nil, fmt.Errorf("mermaid render error: %s", errStr)
		}
	}

	// Get SVG
	svgVal, ok := result["svg"]
	if !ok || svgVal == nil {
		return nil, fmt.Errorf("mermaid render returned no SVG (result: %v)", result)
	}

	svg, ok := svgVal.(string)
	if !ok || svg == "" {
		return nil, fmt.Errorf("mermaid render returned invalid SVG (result: %v)", result)
	}

	return []byte(svg), nil
}

// Close releases browser resources.
// This should be called when the renderer is no longer needed.
func (r *Renderer) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cancel != nil {
		r.cancel()
	}
	if r.allocCancel != nil {
		r.allocCancel()
	}
	if r.tempDir != "" {
		if err := os.RemoveAll(r.tempDir); err != nil {
			log.Error().Err(err).Str("path", r.tempDir).Msg("failed to remove temp directory")
		}
	}
	r.initialized = false
}
