package mermaid

import (
	"bytes"

	"github.com/rs/zerolog/log"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type Extender struct {
	Renderer *Renderer
}

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

type CodeBlock struct {
	ast.BaseBlock
	Code []byte
}

func (n *CodeBlock) Kind() ast.NodeKind {
	return KindMermaidBlock
}

func (n *CodeBlock) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindMermaidBlock = ast.NewNodeKind("MermaidBlock")

type transformer struct{}

func (t *transformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	source := reader.Source()

	var toReplace []*ast.FencedCodeBlock

	err := ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
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
	if err != nil {
		log.Fatal().Err(err).Msg("ast.Walk() failed. Something seriously wrong...")
		return
	}

	for _, fcb := range toReplace {
		var code bytes.Buffer
		lines := fcb.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			code.Write(line.Value(source))
		}

		mermaidBlock := &CodeBlock{
			Code: code.Bytes(),
		}

		parent := fcb.Parent()
		parent.ReplaceChild(parent, fcb, mermaidBlock)
	}
}

type htmlRenderer struct {
	renderer *Renderer
}

func (r *htmlRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindMermaidBlock, r.renderMermaidBlock)
}

func (r *htmlRenderer) renderMermaidBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*CodeBlock)

	if r.renderer == nil {
		_, _ = w.WriteString(`<div class="mermaid">`)
		_, _ = w.Write(n.Code)
		_, _ = w.WriteString("</div>\n")
		return ast.WalkContinue, nil
	}

	svg, err := r.renderer.Render(string(n.Code))
	if err != nil {
		_, _ = w.WriteString("<!-- mermaid render error: ")
		_, _ = w.WriteString(err.Error())
		_, _ = w.WriteString(" -->\n")
		_, _ = w.WriteString(`<pre class="mermaid-error"><code>`)
		_, _ = w.Write(n.Code)
		_, _ = w.WriteString("</code></pre>\n")
		return ast.WalkContinue, nil
	}

	_, _ = w.WriteString(`<div class="mermaid">`)
	_, _ = w.Write(svg)
	_, _ = w.WriteString("</div>\n")

	return ast.WalkContinue, nil
}
