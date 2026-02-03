package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <output-file> <version>\n", os.Args[0])
		os.Exit(1)
	}
	outputPath := os.Args[1]
	katexVersion := os.Args[2]

	katexCSSURL := "https://cdn.jsdelivr.net/npm/katex@" + katexVersion + "/dist/katex.min.css"
	katexFontURL := "https://cdn.jsdelivr.net/npm/katex@" + katexVersion + "/dist/fonts/"

	fmt.Fprintf(os.Stderr, "Fetching KaTeX CSS v%s...\n", katexVersion)
	css, err := fetchURL(katexCSSURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch CSS: %v\n", err)
		os.Exit(1)
	}

	fontURLRegex := regexp.MustCompile(`url\(fonts/([^)]+)\)`)
	matches := fontURLRegex.FindAllStringSubmatch(css, -1)

	fontCache := make(map[string]string)
	for _, match := range matches {
		fontFile := match[1]
		if _, exists := fontCache[fontFile]; exists {
			continue
		}

		fontURL := katexFontURL + fontFile
		fmt.Fprintf(os.Stderr, "Fetching font: %s\n", fontFile)

		fontData, err := fetchURLBytes(fontURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to fetch font %s: %v\n", fontFile, err)
			continue
		}

		mimeType := getMimeType(fontFile)
		dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(fontData))
		fontCache[fontFile] = dataURL
	}

	result := fontURLRegex.ReplaceAllStringFunc(css, func(match string) string {
		submatches := fontURLRegex.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		fontFile := submatches[1]
		if dataURL, exists := fontCache[fontFile]; exists {
			return fmt.Sprintf("url(%s)", dataURL)
		}
		return match
	})

	if err := os.WriteFile(outputPath, []byte(result), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write output: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Written to %s (%d bytes)\n", outputPath, len(result))
}

func fetchURL(url string) (string, error) {
	data, err := fetchURLBytes(url)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func fetchURLBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func getMimeType(filename string) string {
	switch {
	case strings.HasSuffix(filename, ".woff2"):
		return "font/woff2"
	case strings.HasSuffix(filename, ".woff"):
		return "font/woff"
	case strings.HasSuffix(filename, ".ttf"):
		return "font/ttf"
	default:
		return "application/octet-stream"
	}
}
