package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"mdflux/internal/pkg/mdflux/config"
	"mdflux/internal/pkg/mdflux/converter"
	"mdflux/internal/pkg/mdflux/mermaid"
	"mdflux/internal/pkg/mdflux/pdf"
	"mdflux/web"
)

func main() {
	fmt.Fprintf(os.Stderr, "mdflux %s\n", FullVersion())

	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cfg, err := config.LoadAndParse()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration parameters and files")
	}

	if err := setupLogging(cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to setup logging")
	}

	templates, err := converter.ParseTemplates(web.TemplateFS)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse templates")
	}

	log.Debug().Str("config_file", viper.ConfigFileUsed()).Msg("Configuration file used")
	log.Debug().Interface("config", cfg).Msg("Configuration parameters")

	if err := run(cfg, templates); err != nil {
		log.Fatal().Err(err).Msg("Conversion failed")
	}
}

func setupLogging(cfg *config.Config) error {
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Error().Str("log_level", cfg.LogLevel).Msg("Invalid log level, defaulting to info")
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	if cfg.LogFile != "" {
		logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		log.Logger = zerolog.New(logFile).With().Timestamp().Logger()
	}

	return nil
}

func run(cfg *config.Config, templates *converter.Templates) error {
	var input io.Reader

	if cfg.Input == "" || cfg.Input == "-" {
		input = os.Stdin
		log.Debug().Msg("Reading from stdin")
	} else {
		f, err := os.Open(cfg.Input)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Warn().Err(err).Msg("Failed to close input file")
			}
		}()
		input = f
		log.Debug().Str("file", cfg.Input).Msg("Reading from file")
	}

	var mermaidRenderer *mermaid.Renderer
	if cfg.Extensions.Mermaid {
		chromePath := ""
		if cfg.PDF.Chrome.Mode == "manual" {
			chromePath = cfg.PDF.Chrome.Path
		}
		mermaidRenderer = mermaid.NewRenderer(chromePath)
		defer mermaidRenderer.Close()
		log.Debug().Str("chrome_path", chromePath).Msg("Mermaid server-side rendering enabled")
	}

	conv := converter.New(converter.Options{
		Unsafe:              cfg.HTML.Unsafe,
		HardWraps:           cfg.HTML.HardWraps,
		XHTML:               cfg.HTML.XHTML,
		Theme:               cfg.Theme,
		EastAsianLineBreaks: cfg.HTML.EastAsianLineBreaks,
		MermaidRenderer:     mermaidRenderer,
		Extensions: converter.ExtensionOptions{
			Table:          cfg.Extensions.Table,
			Strikethrough:  cfg.Extensions.Strikethrough,
			Linkify:        cfg.Extensions.Linkify,
			TaskList:       cfg.Extensions.TaskList,
			DefinitionList: cfg.Extensions.DefinitionList,
			Footnote:       cfg.Extensions.Footnote,
			Typographer:    cfg.Extensions.Typographer,
			CJK:            cfg.Extensions.CJK,
			D2: converter.D2Options{
				Enabled: cfg.Extensions.D2.Enabled,
				Layout:  cfg.Extensions.D2.Layout,
				ThemeID: cfg.Extensions.D2.ThemeID,
			},
			KaTeX:   cfg.Extensions.KaTeX,
			Mermaid: cfg.Extensions.Mermaid,
		},
	}, templates)

	format := cfg.Format
	if format == "" {
		format = "html"
	}

	log.Info().Str("format", format).Msg("Starting conversion")

	if format == "pdf" {
		return runPDFConversion(cfg, conv, input)
	}

	return runHTMLConversion(cfg, conv, input)
}

func runHTMLConversion(cfg *config.Config, conv *converter.Converter, input io.Reader) error {
	var output io.Writer

	if cfg.Output == "" || cfg.Output == "-" {
		output = os.Stdout
		log.Debug().Msg("Writing to stdout")
	} else {
		f, err := os.Create(cfg.Output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Warn().Err(err).Msg("Failed to close output file")
			}
		}()
		output = f
		log.Debug().Str("file", cfg.Output).Msg("Writing to file")
	}

	if err := conv.ConvertReader(input, output); err != nil {
		return fmt.Errorf("conversion error: %w", err)
	}

	log.Info().Msg("Conversion completed successfully")
	return nil
}

func runPDFConversion(cfg *config.Config, conv *converter.Converter, input io.Reader) error {
	if cfg.Output == "" || cfg.Output == "-" {
		return fmt.Errorf("PDF output requires a file path, cannot write to stdout")
	}

	tmpFile, err := os.CreateTemp("", "mdflux-*.html")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tempHTMLPath := tmpFile.Name()
	defer func() {
		if err := os.Remove(tempHTMLPath); err != nil {
			log.Warn().Err(err).Msg("Failed to remove temporary file")
		}
	}()

	log.Debug().Str("temp_file", tempHTMLPath).Msg("Created temporary HTML file")

	if err := conv.ConvertReader(input, tmpFile); err != nil {
		if closeErr := tmpFile.Close(); closeErr != nil {
			log.Warn().Err(closeErr).Msg("Failed to close temporary file")
		}
		return fmt.Errorf("conversion error: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		log.Warn().Err(err).Msg("Failed to close temporary file")
	}

	absHTMLPath, err := filepath.Abs(tempHTMLPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	absPDFPath, err := filepath.Abs(cfg.Output)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for output: %w", err)
	}

	pdfOpts := pdf.Options{
		PageSize:     cfg.PDF.PageSize,
		Landscape:    cfg.PDF.Landscape,
		Scale:        cfg.PDF.Scale,
		MarginTop:    cfg.PDF.MarginTop,
		MarginBottom: cfg.PDF.MarginBottom,
		MarginLeft:   cfg.PDF.MarginLeft,
		MarginRight:  cfg.PDF.MarginRight,
		ChromeMode:   cfg.PDF.Chrome.Mode,
		ChromePath:   cfg.PDF.Chrome.Path,
	}

	log.Debug().Str("html_path", absHTMLPath).Str("pdf_path", absPDFPath).Msg("Rendering PDF")

	if err := pdf.RenderHTMLToPDF(absHTMLPath, absPDFPath, pdfOpts); err != nil {
		return fmt.Errorf("PDF rendering failed: %w", err)
	}

	log.Info().Str("output", absPDFPath).Msg("PDF conversion completed successfully")
	return nil
}
