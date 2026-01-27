# mdflux

**mdflux** is a high-performance, minimalist CLI utility and Go library designed to convert Commonmark Markdown into production-ready HTML or PDF documents. Unlike heavy Node.js alternatives, this tool is compiled to a single binary with dependency on a webkit browser like Chrome or Edge for PDF generation.

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

All options can be set via environment variables with the `MDFLUX_` prefix:

```bash
export MDFLUX_THEME=dark
export MDFLUX_FORMAT=pdf
```

### Config File

Create a `mdflux.cfg.toml` file:

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

### Extension Options

All extensions are enabled by default. Set to `false` to disable.

| Option | Default | Description |
| --- | --- | --- |
| `table` | `true` | GFM table syntax with headers, alignment, and cell formatting. |
| `strikethrough` | `true` | Strike-through text using `~~deleted~~` syntax. |
| `linkify` | `true` | Auto-detect and convert URLs and email addresses to clickable links. |
| `task_list` | `true` | GitHub-style task lists with `- [x]` and `- [ ]` checkbox syntax. |
| `definition_list` | `true` | Definition lists using term/definition pairs (PHP Markdown Extra syntax). |
| `footnote` | `true` | Footnote references `[^1]` with auto-generated footnotes section. |
| `typographer` | `true` | Smart typography: straight quotes → curly quotes, `--` → en-dash, `---` → em-dash, `...` → ellipsis. |
| `cjk` | `true` | Optimized rendering for Chinese, Japanese, and Korean text. |
| `katex` | `true` | LaTeX math rendering. Inline: `$E=mc^2$`. Display: `$$\int_0^\infty$$`. |
| `mermaid` | `true` | Mermaid diagrams in fenced code blocks with `mermaid` language identifier. |

### D2 Diagram Options

D2 is a declarative diagramming language. Configure under `[extensions.d2]`:

| Option | Default | Description |
| --- | --- | --- |
| `enabled` | `true` | Enable D2 diagram rendering in fenced code blocks with `d2` language identifier. |
| `layout` | `dagre` | Layout engine. `dagre` for directed graphs, `elk` for more complex layouts. |
| `theme_id` | `0` | D2 theme ID. `0` is default, other values apply different color schemes. |

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
