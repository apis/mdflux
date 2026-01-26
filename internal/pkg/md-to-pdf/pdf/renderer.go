package pdf

import (
	"context"
	"fmt"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Options struct {
	PageSize     string
	Landscape    bool
	Scale        float64
	MarginTop    float64
	MarginBottom float64
	MarginLeft   float64
	MarginRight  float64
	ChromeMode   string
	ChromePath   string
}

// DefaultOptions returns Options with sensible defaults for PDF rendering.
func DefaultOptions() Options {
	return Options{
		PageSize:     "A4",
		Landscape:    false,
		Scale:        0.8,
		MarginTop:    0.5,
		MarginBottom: 0.5,
		MarginLeft:   0.5,
		MarginRight:  0.5,
		ChromeMode:   "auto",
		ChromePath:   "",
	}
}

func RenderHTMLToPDF(htmlFilePath, pdfFilePath string, opts Options) error {
	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.NoSandbox,
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("no-zygote", true),
		chromedp.Flag("disable-extensions", true),
	)

	if opts.ChromeMode == "manual" && opts.ChromePath != "" {
		allocOpts = append(allocOpts, chromedp.ExecPath(opts.ChromePath))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	defer allocCancel()

	ctx, ctxCancel := chromedp.NewContext(allocCtx)
	defer ctxCancel()

	var buf []byte
	if err := chromedp.Run(ctx, printToPDFTasks(htmlFilePath, &buf, opts)); err != nil {
		return fmt.Errorf("chromedp.Run() failed: %w", err)
	}

	if err := os.WriteFile(pdfFilePath, buf, 0644); err != nil {
		return fmt.Errorf("os.WriteFile() failed: %w", err)
	}

	return nil
}

func printToPDFTasks(htmlPath string, buffer *[]byte, opts Options) chromedp.Tasks {
	fileURL := "file://" + htmlPath

	paperWidth, paperHeight := getPaperSize(opts.PageSize)
	if opts.Landscape {
		paperWidth, paperHeight = paperHeight, paperWidth
	}

	scale := opts.Scale
	if scale <= 0 {
		scale = 0.8
	}

	return chromedp.Tasks{
		chromedp.Navigate(fileURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			pdfConfig := page.PrintToPDF().
				WithPrintBackground(true).
				WithScale(scale).
				WithPaperWidth(paperWidth).
				WithPaperHeight(paperHeight).
				WithMarginTop(opts.MarginTop).
				WithMarginBottom(opts.MarginBottom).
				WithMarginLeft(opts.MarginLeft).
				WithMarginRight(opts.MarginRight)

			pdfBuffer, _, err := pdfConfig.Do(ctx)
			if err != nil {
				return err
			}
			*buffer = pdfBuffer
			return nil
		}),
	}
}

func getPaperSize(size string) (width, height float64) {
	switch size {
	case "Letter":
		return 8.5, 11.0
	case "Legal":
		return 8.5, 14.0
	case "A4":
		fallthrough
	default:
		return 8.27, 11.69
	}
}
