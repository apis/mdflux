# mdflux

> Transform your Markdown into high-fidelity HTML and PDF with zero friction.

**mdflux** is a high-performance, minimalist CLI utility and Go library designed to convert Commonmark Markdown into production-ready HTML or PDF documents. Unlike heavy Node.js alternatives, this tool is compiled to a single binary with dependency on a webkit browser like Chrome or Edge for PDF generation.

## Demo

See mdflux in action with Mermaid diagrams, D2 architecture diagrams, and KaTeX mathematical expressions—all rendered server-side to embedded SVG with no client-side JavaScript required:

| Format | Link |
| --- | --- |
| Markdown | [docs/demo.md](docs/demo.md) |
| HTML | [docs/demo.html](docs/demo.html) |
| PDF | [docs/demo.pdf](docs/demo.pdf) |

---

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

Download the pre-compiled binary from the [Releases](https://github.com/apis/mdflux/releases) page or install via Go:

```bash
go install github.com/apis/mdflux@latest
```

---

## Usage

### Command Line Interface

Convert a markdown file to PDF:

```bash
mdflux -i input.md -o output.pdf -f pdf
```

Convert to HTML with dark theme:

```bash
mdflux -i input.md -o output.html -t dark
```

Read from stdin and output HTML to stdout:

```bash
cat input.md | mdflux
```

---

## Configuration

Configuration is loaded from multiple sources with the following precedence (highest to lowest):

1. Command-line flags
2. Environment variables
3. Config file
4. Default values

### CLI Flags

| Flag | Shorthand | Description | Default |
| --- | --- | --- | --- |
| `--config` | `-c` | Path to config file | (auto-detect) |
| `--input` | `-i` | Input markdown file (use `-` for stdin) | stdin |
| `--output` | `-o` | Output file (use `-` for stdout) | stdout |
| `--format` | `-f` | Output format (`html`, `pdf`) | `html` |
| `--theme` | `-t` | Color theme (`auto`, `light`, `dark`) | `auto` |
| `--log_level` | `-l` | Log level (`debug`, `info`, `warn`, `error`) | `info` |
| `--log_file` | | Log file path | stderr |
| `--help` | `-?` | Display help | |

### Environment Variables

All options can be set via environment variables with the `MDFLUX_` prefix:

```bash
export MDFLUX_INPUT=input.md
export MDFLUX_OUTPUT=output.pdf
export MDFLUX_FORMAT=pdf
export MDFLUX_THEME=dark
export MDFLUX_LOG_LEVEL=debug
```

### Config File

mdflux searches for `mdflux.cfg.toml` in these locations (in order):

1. Current working directory (`./`)
2. User config directory (`$HOME/.config/mdflux/`)
3. System config directory (`/etc/mdflux/`)

Or specify a custom path with `-c /path/to/config.toml`.

Example `mdflux.cfg.toml`:

```toml
input = ""
output = ""
format = "html"
theme = "auto"
log_level = "info"
log_file = ""

[html]
unsafe = false
hard_wraps = false
xhtml = false
east_asian_line_breaks = "simple"

[pdf]
page_size = "A4"
landscape = false
scale = 0.8
margin_top = 0.5
margin_bottom = 0.5
margin_left = 0.5
margin_right = 0.5

[pdf.chrome]
mode = "auto"
path = ""

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

---

## HTML Options

| Option | Default | Description |
| --- | --- | --- |
| `theme` | `auto` | Color theme. `auto` follows system preference, `light` or `dark` for fixed themes. |
| `unsafe` | `false` | Allow raw HTML in markdown. When `false`, HTML tags are escaped. Enable for trusted content only. |
| `hard_wraps` | `false` | Render single line breaks as `<br>`. When `false`, single newlines become spaces (standard Commonmark). |
| `xhtml` | `false` | Output XHTML 1.0 Strict instead of HTML5. Produces self-closing tags (`<br />`) and XML declaration. |
| `east_asian_line_breaks` | `simple` | Line break handling for CJK text. `simple` removes breaks between wide characters. `css3draft` follows CSS Text Level 3 rules. |

---

## PDF Options

| Option | Default | Description |
| --- | --- | --- |
| `page_size` | `A4` | Page size (`A4`, `Letter`, `Legal`). |
| `landscape` | `false` | Use landscape orientation. |
| `scale` | `0.8` | Scale factor for rendering (`0.1` - `2.0`). |
| `margin_top` | `0.5` | Top margin in inches. |
| `margin_bottom` | `0.5` | Bottom margin in inches. |
| `margin_left` | `0.5` | Left margin in inches. |
| `margin_right` | `0.5` | Right margin in inches. |

### Chrome Configuration

PDF rendering uses headless Chrome/Chromium. Configure under `[pdf.chrome]`:

| Option | Default | Description |
| --- | --- | --- |
| `mode` | `auto` | Chrome detection mode. `auto` finds Chrome automatically, `manual` uses specified path. |
| `path` | `""` | Path to Chrome/Chromium executable (only used when `mode = "manual"`). |

---

## Extension Options

All extensions are enabled by default. Set to `false` to disable.

| Option | Default | Description |
| --- | --- | --- |
| `table` | `true` | GFM table syntax with headers, alignment, and cell formatting. |
| `strikethrough` | `true` | Strike-through text using `~~deleted~~` syntax. |
| `linkify` | `true` | Auto-detect and convert URLs and email addresses to clickable links. |
| `task_list` | `true` | GitHub-style task lists with `- [x]` and `- [ ]` checkbox syntax. |
| `definition_list` | `true` | Definition lists using term/definition pairs (PHP Markdown Extra syntax). |
| `footnote` | `true` | Footnote references `[^1]` with auto-generated footnotes section. |
| `typographer` | `true` | Smart typography: straight quotes to curly quotes, `--` to en-dash, `---` to em-dash, `...` to ellipsis. |
| `cjk` | `true` | Optimized rendering for Chinese, Japanese, and Korean text. |
| `katex` | `true` | LaTeX math rendering. Inline: `$E=mc^2$`. Display: `$$\int_0^\infty$$`. |
| `mermaid` | `true` | Mermaid diagrams in fenced code blocks with `mermaid` language identifier. Server-side rendered to SVG. |

### D2 Diagram Options

D2 is a declarative diagramming language. Configure under `[extensions.d2]`:

| Option | Default | Description |
| --- | --- | --- |
| `enabled` | `true` | Enable D2 diagram rendering in fenced code blocks with `d2` language identifier. |
| `layout` | `dagre` | Layout engine. `dagre` for directed graphs, `elk` for more complex layouts. |
| `theme_id` | `0` | D2 theme ID. `0` is default, other values apply different color schemes. |

---

## License

Distributed under the MIT License. See `LICENSE` for more information.
