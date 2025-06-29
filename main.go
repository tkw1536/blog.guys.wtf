package main

import (
	"cmp"
	"context"
	"fmt"
	"generator/generator"
	"html/template"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"time"

	_ "embed"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
)

var globals = map[string]any{
	"URL":       "https://blog.guys.wtf",
	"BlogTitle": "High on Code!",
}

//go:embed "templates/index.html"
var indexHTML string
var indexTemplate = mustTemplate(indexHTML, "index.html")

//go:embed templates/list.html
var listHTML string
var listTemplate = mustTemplate(listHTML, "list.html")

var g = generator.Generator{
	Inputs: []generator.Scanner{
		generator.NewStaticScanner("static", []string{"_", "."}),
		generator.NewMarkdownScanner(
			"content",
			func(path string, Metadata map[string]any) bool {
				return Metadata["draft"] != true
			},
			goldmark.WithExtensions(
				extension.GFM,
				highlighting.NewHighlighting(
					highlighting.WithStyle("monokai"),
					highlighting.WithFormatOptions(
						chromahtml.WithLineNumbers(true),
					),
				),
			),
			goldmark.WithRendererOptions(
			//	html.WithUnsafe(),
			)),
	},

	Indexes: []generator.IndexTemplate{
		{
			Path:     "index.html",
			Template: listTemplate,
			CompareFunc: func(left, right generator.IndexEntry) int {
				lDate, _ := left.Metadata["date"].(string)
				rDate, _ := right.Metadata["date"].(string)
				return cmp.Or(
					strings.Compare(rDate, lDate), // descending by date
					strings.Compare(left.Path, right.Path),
				)
			},
		},
	},

	ContentTemplate: generator.ContentTemplate{
		Template: indexTemplate,
		Globals:  globals,
	},

	PostProcessors: []generator.PostProcessor{
		generator.MinifyPostProcessor,
	},

	Output: generator.NewNativeFileWriter("public", true),
}

func main() {
	// when done, exit with the code!
	var exitCode int
	defer func() {
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// create a new logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	// and run
	if err := g.Run(ctx, logger); err != nil {
		exitCode = 1
	}
}

func mustTemplate(src, name string) *template.Template {
	return template.Must(template.New(name).Funcs(templateFuncs).Parse(src))
}

var templateFuncs = template.FuncMap{
	"date": func(arg string) (string, error) {
		date, err := time.Parse("2006-01-02", arg)
		if err != nil {
			return "", fmt.Errorf("failed to parse date %q: %w", arg, err)
		}

		day := date.Day()

		return fmt.Sprintf(
			"%s %d%s %d",
			date.Format("January"),
			day, getSuffix(day),
			date.Year(),
		), nil
	},
}

func getSuffix(day int) string {
	if day >= 11 && day <= 13 {
		return "th"
	}
	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}
