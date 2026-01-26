package converter

import (
	"fmt"
	"io"

	d2 "github.com/FurqanSoftware/goldmark-d2"
	"github.com/FurqanSoftware/goldmark-katex"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	"md-to-pdf/internal/pkg/md-to-pdf/mermaid"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2layouts/d2elklayout"
)

type Options struct {
	Unsafe              bool
	HardWraps           bool
	XHTML               bool
	Theme               string
	EastAsianLineBreaks string
	Extensions          ExtensionOptions
	MermaidRenderer     *mermaid.Renderer
}

type ExtensionOptions struct {
	Table          bool
	Strikethrough  bool
	Linkify        bool
	TaskList       bool
	DefinitionList bool
	Footnote       bool
	Typographer    bool
	CJK            bool
	D2             D2Options
	KaTeX          bool
	Mermaid        bool
}

type D2Options struct {
	Enabled bool
	Layout  string
	ThemeID int64
}

type Converter struct {
	markdown   goldmark.Markdown
	templates  *Templates
	xhtml      bool
	theme      string
	extensions ExtensionOptions
}

func New(opts Options, templates *Templates) *Converter {
	var htmlOpts []html.Option

	if opts.Unsafe {
		htmlOpts = append(htmlOpts, html.WithUnsafe())
	}
	if opts.HardWraps {
		htmlOpts = append(htmlOpts, html.WithHardWraps())
	}
	if opts.XHTML {
		htmlOpts = append(htmlOpts, html.WithXHTML())
	}
	switch opts.EastAsianLineBreaks {
	case "simple":
		htmlOpts = append(htmlOpts, html.WithEastAsianLineBreaks(html.EastAsianLineBreaksSimple))
	case "css3draft":
		htmlOpts = append(htmlOpts, html.WithEastAsianLineBreaks(html.EastAsianLineBreaksCSS3Draft))
	}

	htmlRenderer := html.NewRenderer(htmlOpts...)

	var gmOpts []goldmark.Option

	gmOpts = append(gmOpts, goldmark.WithRenderer(
		renderer.NewRenderer(
			renderer.WithNodeRenderers(
				util.Prioritized(htmlRenderer, 1000),
			),
		),
	))

	if opts.Extensions.Table {
		gmOpts = append(gmOpts, goldmark.WithExtensions(extension.Table))
	}
	if opts.Extensions.Strikethrough {
		gmOpts = append(gmOpts, goldmark.WithExtensions(extension.Strikethrough))
	}
	if opts.Extensions.Linkify {
		gmOpts = append(gmOpts, goldmark.WithExtensions(extension.Linkify))
	}
	if opts.Extensions.TaskList {
		gmOpts = append(gmOpts, goldmark.WithExtensions(extension.TaskList))
	}
	if opts.Extensions.DefinitionList {
		gmOpts = append(gmOpts, goldmark.WithExtensions(extension.DefinitionList))
	}
	if opts.Extensions.Footnote {
		gmOpts = append(gmOpts, goldmark.WithExtensions(extension.Footnote))
	}
	if opts.Extensions.Typographer {
		gmOpts = append(gmOpts, goldmark.WithExtensions(extension.Typographer))
	}
	if opts.Extensions.CJK {
		gmOpts = append(gmOpts, goldmark.WithExtensions(extension.CJK))
	}

	if opts.Extensions.D2.Enabled {
		var layoutFunc d2graph.LayoutGraph
		switch opts.Extensions.D2.Layout {
		case "elk":
			layoutFunc = d2elklayout.DefaultLayout
		default:
			layoutFunc = d2dagrelayout.DefaultLayout
		}
		themeID := opts.Extensions.D2.ThemeID
		gmOpts = append(gmOpts, goldmark.WithExtensions(&d2.Extender{
			Layout:  layoutFunc,
			ThemeID: &themeID,
		}))
	}

	if opts.Extensions.KaTeX {
		gmOpts = append(gmOpts, goldmark.WithExtensions(&katex.Extender{}))
	}

	if opts.Extensions.Mermaid {
		gmOpts = append(gmOpts, goldmark.WithExtensions(&mermaid.Extender{
			Renderer: opts.MermaidRenderer,
		}))
	}

	md := goldmark.New(gmOpts...)

	return &Converter{
		markdown:   md,
		templates:  templates,
		xhtml:      opts.XHTML,
		theme:      opts.Theme,
		extensions: opts.Extensions,
	}
}

func (c *Converter) Convert(source []byte, w io.Writer) error {
	headerTemplate := "html5-header"
	footerTemplate := "html5-footer"
	if c.xhtml {
		headerTemplate = "xhtml-header"
		footerTemplate = "xhtml-footer"
	}

	data := HeaderData{
		Title:  "Document",
		Styles: c.templates.Styles(),
		Theme:  c.theme,
	}

	if err := c.templates.Template().ExecuteTemplate(w, headerTemplate, data); err != nil {
		return fmt.Errorf("failed to execute header template: %w", err)
	}

	if err := c.markdown.Convert(source, w); err != nil {
		return fmt.Errorf("goldmark conversion failed: %w", err)
	}

	if err := c.templates.Template().ExecuteTemplate(w, footerTemplate, nil); err != nil {
		return fmt.Errorf("failed to execute footer template: %w", err)
	}

	return nil
}

func (c *Converter) ConvertReader(r io.Reader, w io.Writer) error {
	source, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	return c.Convert(source, w)
}
