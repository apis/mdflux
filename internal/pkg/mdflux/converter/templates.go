package converter

import (
	"fmt"
	"io/fs"
	"strings"
	"text/template"

	"mdflux/web"
)

const (
	templatesDir = "templates"
	templateExt  = ".gohtml"
	stylesFile   = "templates/styles.css"
)

type Templates struct {
	templates *template.Template
	styles    string
}

type HeaderData struct {
	Title  string
	Styles string
	Theme  string
}

func ParseTemplates(templateFS fs.FS) (*Templates, error) {
	styles, err := fs.ReadFile(templateFS, stylesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read styles.css: %w", err)
	}

	tmpl, err := parseTemplatesRecursive(templateFS, templatesDir, templateExt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	combinedStyles := web.KaTeXCSS + "\n" + string(styles)

	return &Templates{
		templates: tmpl,
		styles:    combinedStyles,
	}, nil
}

func parseTemplatesRecursive(templateFS fs.FS, dir string, ext string) (*template.Template, error) {
	root := template.New("")

	err := fs.WalkDir(templateFS, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ext) {
			return nil
		}

		content, err := fs.ReadFile(templateFS, path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}

		_, err = root.Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return root, nil
}

func (t *Templates) Styles() string {
	return t.styles
}

func (t *Templates) Template() *template.Template {
	return t.templates
}
