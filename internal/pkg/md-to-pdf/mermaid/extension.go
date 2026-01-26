package mermaid

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Extender is a goldmark extension that renders mermaid code blocks to SVG.
type Extender struct {
	Renderer *Renderer
}

// Extend implements goldmark.Extender.
func (e *Extender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&transformer{}, 100),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&htmlRenderer{renderer: e.Renderer}, 100),
		),
	)
}

// MermaidBlock is a custom AST node for mermaid diagrams.
type MermaidBlock struct {
	ast.BaseBlock
	Code []byte
}

// Kind returns the kind of the node.
func (n *MermaidBlock) Kind() ast.NodeKind {
	return KindMermaidBlock
}

// Dump dumps the node to stdout for debugging.
func (n *MermaidBlock) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// KindMermaidBlock is the kind of MermaidBlock.
var KindMermaidBlock = ast.NewNodeKind("MermaidBlock")

// transformer converts mermaid fenced code blocks to MermaidBlock nodes.
type transformer struct{}

// Transform implements parser.ASTTransformer.
func (t *transformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	source := reader.Source()

	var toReplace []*ast.FencedCodeBlock

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		fcb, ok := n.(*ast.FencedCodeBlock)
		if !ok {
			return ast.WalkContinue, nil
		}

		lang := fcb.Language(source)
		if !bytes.Equal(lang, []byte("mermaid")) {
			return ast.WalkContinue, nil
		}

		toReplace = append(toReplace, fcb)
		return ast.WalkContinue, nil
	})

	for _, fcb := range toReplace {
		var code bytes.Buffer
		lines := fcb.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			code.Write(line.Value(source))
		}

		mermaidBlock := &MermaidBlock{
			Code: code.Bytes(),
		}

		parent := fcb.Parent()
		parent.ReplaceChild(parent, fcb, mermaidBlock)
	}
}

// htmlRenderer renders MermaidBlock nodes to HTML.
type htmlRenderer struct {
	renderer *Renderer
}

// RegisterFuncs implements renderer.NodeRenderer.
func (r *htmlRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindMermaidBlock, r.renderMermaidBlock)
}

func (r *htmlRenderer) renderMermaidBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*MermaidBlock)

	if r.renderer == nil {
		// Fallback: render as client-side mermaid div
		_, _ = w.WriteString(`<div class="mermaid">`)
		_, _ = w.Write(n.Code)
		_, _ = w.WriteString("</div>\n")
		return ast.WalkContinue, nil
	}

	svg, err := r.renderer.Render(string(n.Code))
	if err != nil {
		// On error, output the error as a comment and the original code
		_, _ = w.WriteString("<!-- mermaid render error: ")
		_, _ = w.WriteString(err.Error())
		_, _ = w.WriteString(" -->\n")
		_, _ = w.WriteString(`<pre class="mermaid-error"><code>`)
		_, _ = w.Write(n.Code)
		_, _ = w.WriteString("</code></pre>\n")
		return ast.WalkContinue, nil
	}

	// Wrap SVG in a div for styling consistency
	_, _ = w.WriteString(`<div class="mermaid">`)
	_, _ = w.Write(svg)
	_, _ = w.WriteString("</div>\n")

	return ast.WalkContinue, nil
}
