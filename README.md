# md-to-pdf

**md-to-pdf** is a high-performance, minimalist CLI utility and Go library designed to convert Commonmark Markdown into production-ready HTML or PDF documents. Unlike heavy Node.js alternatives, this tool is compiled to a single binary with dependency on a webkit browser like Chrome or Edge for PDF generation.

## Features

* **Performance:** Built with Go for near-instantaneous conversion of large documents
* **Multiple Output Formats:** Generate HTML5, PDF, or XHTML 1.0 Strict output
* **Theming:** Built-in auto, light, and dark themes with CSS variable customization
* **Flexible Input:** Read from files or stdin
* **Configuration:** TOML config files, environment variables, and CLI flags

## Markdown Extensions

All extensions are enabled by default and can be toggled via configuration:

| Extension | Description |
| --- | --- |
| **Tables** | GFM (GitHub Flavored Markdown) table syntax |
| **Strikethrough** | `~~text~~` syntax for struck-through text |
| **Linkify** | Auto-detection and linking of URLs and email addresses |
| **Task Lists** | GitHub-style task lists: `- [x] done` and `- [ ] todo` |
| **Definition Lists** | PHP Markdown Extra definition list syntax |
| **Footnotes** | `[^1]` footnote reference syntax with automatic section generation |
| **Typographer** | Smart quotes, dashes, and ellipses |
| **CJK** | Optimized text handling for Chinese, Japanese, and Korean |
| **KaTeX** | LaTeX math rendering for inline (`$...$`) and display (`$$...$$`) equations |
| **Mermaid** | Flowcharts, sequence diagrams, class diagrams, Gantt charts, and more |
| **D2** | Declarative diagrams with multiple layout engines (Dagre, ELK) |

---

## Installation

### CLI Tool

Download the pre-compiled binary from the [Releases](https://github.com/apis/md-to-pdf/releases) page or install via Go:

```bash
go install github.com/apis/md-to-pdf@latest
```

---

## Usage

### Command Line Interface

Convert a markdown file to PDF:

```bash
md-to-pdf -i input.md -o output.pdf -f pdf
```

Convert to HTML with dark theme:

```bash
md-to-pdf -i input.md -o output.html -t dark
```

Read from stdin and output HTML to stdout:

```bash
cat input.md | md-to-pdf
```

---

## Configuration

### CLI Flags

| Flag | Shorthand | Description | Default |
| --- | --- | --- | --- |
| `--config` | `-c` | Path to config file | (auto-detect) |
| `--input` | `-i` | Input markdown file | stdin |
| `--output` | `-o` | Output file path | stdout (HTML) |
| `--format` | `-f` | Output format (`html`, `pdf`) | `html` |
| `--theme` | `-t` | Color theme (`auto`, `light`, `dark`) | `auto` |
| `--log-level` | `-l` | Log level (`debug`, `info`, `warn`, `error`) | `info` |
| `--log-file` |  | Log file path | stderr |
| `--help` | `-?` | Display help |  |

### Environment Variables

All options can be set via environment variables with the `MD_TO_PDF_` prefix:

```bash
export MD_TO_PDF_THEME=dark
export MD_TO_PDF_FORMAT=pdf
```

### Config File

Create a `md-to-pdf.cfg.toml` file:

```toml
[html]
theme = "auto"
unsafe = false
hard_wraps = false
xhtml = false
east_asian_line_breaks = "simple"

[pdf]
page_size = "A4"
landscape = false
scale = 0.8
margin_top = 0.4
margin_bottom = 0.4
margin_left = 0.4
margin_right = 0.4

[extensions]
table = true
strikethrough = true
linkify = true
task_list = true
definition_list = true
footnote = true
typographer = true
cjk = true
katex = true
mermaid = true

[extensions.d2]
enabled = true
layout = "dagre"
theme_id = 0
```

### HTML Options

| Option | Default | Description |
| --- | --- | --- |
| `unsafe` | `false` | Allow raw HTML in markdown. When `false`, HTML tags are escaped. Enable for trusted content only. |
| `hard_wraps` | `false` | Render single line breaks as `<br>`. When `false`, single newlines become spaces (standard Commonmark). |
| `xhtml` | `false` | Output XHTML 1.0 Strict instead of HTML5. Produces self-closing tags (`<br />`) and XML declaration. |
| `east_asian_line_breaks` | `simple` | Line break handling for CJK text. `simple` removes breaks between wide characters. `css3draft` follows CSS Text Level 3 rules. |

---

## PDF Options

| Option | Values | Default |
| --- | --- | --- |
| `page_size` | `A4`, `Letter`, `Legal` | `A4` |
| `landscape` | `true`, `false` | `false` |
| `scale` | `0.1` - `2.0` | `0.8` |
| `margin_*` | inches | `0.4` |

---

## License

Distributed under the MIT License. See `LICENSE` for more information.
