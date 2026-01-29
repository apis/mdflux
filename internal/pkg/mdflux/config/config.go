package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	envPrefix = "MDFLUX"

	helpKey     = "help"
	configKey   = "config"
	inputKey    = "input"
	outputKey   = "output"
	formatKey   = "format"
	logLevelKey = "log_level"
	logFileKey  = "log_file"
	themeKey    = "theme"

	defaultLogLevel      = "info"
	defaultTheme         = "auto"
	defaultFormat        = "html"
	defaultD2Layout      = "dagre"
	defaultD2ThemeID     = int64(0)
	defaultPDFPageSize   = "A4"
	defaultPDFScale      = 0.8
	defaultPDFMargin     = 0.5
	defaultPDFChromeMode = "auto"
)

type Config struct {
	Input      string           `mapstructure:"input"`
	Output     string           `mapstructure:"output"`
	Format     string           `mapstructure:"format"`
	Theme      string           `mapstructure:"theme"`
	LogLevel   string           `mapstructure:"log_level"`
	LogFile    string           `mapstructure:"log_file"`
	HTML       HTMLConfig       `mapstructure:"html"`
	PDF        PDFConfig        `mapstructure:"pdf"`
	Extensions ExtensionsConfig `mapstructure:"extensions"`
}

type PDFConfig struct {
	PageSize     string       `mapstructure:"page_size"`
	Landscape    bool         `mapstructure:"landscape"`
	Scale        float64      `mapstructure:"scale"`
	MarginTop    float64      `mapstructure:"margin_top"`
	MarginBottom float64      `mapstructure:"margin_bottom"`
	MarginLeft   float64      `mapstructure:"margin_left"`
	MarginRight  float64      `mapstructure:"margin_right"`
	Chrome       ChromeConfig `mapstructure:"chrome"`
}

type ChromeConfig struct {
	Mode string `mapstructure:"mode"`
	Path string `mapstructure:"path"`
}

type HTMLConfig struct {
	Unsafe              bool   `mapstructure:"unsafe"`
	HardWraps           bool   `mapstructure:"hard_wraps"`
	XHTML               bool   `mapstructure:"xhtml"`
	EastAsianLineBreaks string `mapstructure:"east_asian_line_breaks"`
}

type ExtensionsConfig struct {
	Table          bool     `mapstructure:"table"`
	Strikethrough  bool     `mapstructure:"strikethrough"`
	Linkify        bool     `mapstructure:"linkify"`
	TaskList       bool     `mapstructure:"task_list"`
	DefinitionList bool     `mapstructure:"definition_list"`
	Footnote       bool     `mapstructure:"footnote"`
	Typographer    bool     `mapstructure:"typographer"`
	CJK            bool     `mapstructure:"cjk"`
	D2             D2Config `mapstructure:"d2"`
	KaTeX          bool     `mapstructure:"katex"`
	Mermaid        bool     `mapstructure:"mermaid"`
}

type D2Config struct {
	Enabled bool   `mapstructure:"enabled"`
	Layout  string `mapstructure:"layout"`
	ThemeID int64  `mapstructure:"theme_id"`
}

func LoadAndParse() (*Config, error) {
	viper.SetDefault(logLevelKey, defaultLogLevel)
	viper.SetDefault(themeKey, defaultTheme)
	viper.SetDefault(formatKey, defaultFormat)
	viper.SetDefault("extensions.d2.layout", defaultD2Layout)
	viper.SetDefault("extensions.d2.theme_id", defaultD2ThemeID)
	viper.SetDefault("pdf.page_size", defaultPDFPageSize)
	viper.SetDefault("pdf.scale", defaultPDFScale)
	viper.SetDefault("pdf.margin_top", defaultPDFMargin)
	viper.SetDefault("pdf.margin_bottom", defaultPDFMargin)
	viper.SetDefault("pdf.margin_left", defaultPDFMargin)
	viper.SetDefault("pdf.margin_right", defaultPDFMargin)
	viper.SetDefault("pdf.chrome.mode", defaultPDFChromeMode)

	viper.SetDefault("extensions.table", true)
	viper.SetDefault("extensions.strikethrough", true)
	viper.SetDefault("extensions.linkify", true)
	viper.SetDefault("extensions.task_list", true)
	viper.SetDefault("extensions.typographer", true)
	viper.SetDefault("extensions.katex", true)
	viper.SetDefault("extensions.mermaid", true)
	viper.SetDefault("extensions.d2.enabled", true)

	flagSet := pflag.NewFlagSet("mdflux", pflag.ContinueOnError)
	flagSet.Usage = func() {}

	help := flagSet.BoolP(helpKey, "?", false, "Display help information")
	flagSet.StringP(configKey, "c", "", "Path to config file")
	flagSet.StringP(inputKey, "i", "", "Input markdown file (use - for stdin)")
	flagSet.StringP(outputKey, "o", "", "Output file (use - for stdout)")
	flagSet.StringP(formatKey, "f", defaultFormat, "Output format (html, pdf)")
	flagSet.StringP(logLevelKey, "l", defaultLogLevel, "Log level (debug, info, warn, error)")
	flagSet.String(logFileKey, "", "Log file path")
	flagSet.StringP(themeKey, "t", defaultTheme, "Color theme (auto, light, dark)")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			fmt.Println(flagSet.FlagUsages())
			os.Exit(0)
		}
		return nil, fmt.Errorf("flagSet.Parse() failed: %w", err)
	}

	if *help {
		fmt.Println("mdflux - Transform your Markdown into high-fidelity HTML and PDF with zero friction.")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  mdflux [flags]")
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println(flagSet.FlagUsages())
		os.Exit(0)
	}

	if err := viper.BindPFlags(flagSet); err != nil {
		return nil, fmt.Errorf("viper.BindPFlags() failed: %w", err)
	}

	configFile, _ := flagSet.GetString(configKey)
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("mdflux.cfg")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.config/mdflux")
		viper.AddConfigPath("/etc/mdflux")
	}

	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	for _, key := range []string{inputKey, outputKey, formatKey, logLevelKey, logFileKey, themeKey} {
		_ = viper.BindEnv(key)
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("viper.ReadInConfig() failed: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("viper.Unmarshal() failed: %w", err)
	}

	return &cfg, nil
}
