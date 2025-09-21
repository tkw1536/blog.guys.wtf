//spellchecker:words main
package main

//spellchecker:words context generator html template slog http signal strings time embed github alecthomas chroma formatters chromahtml yuin goldmark highlighting extension
import (
	"cmp"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"go.tkw01536.de/blog/generator"
	"go.tkw01536.de/blog/generator/output"
	"go.tkw01536.de/blog/generator/scanner"

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
	Inputs: []scanner.Scanner{
		scanner.Static("static", func(name string) bool {
			return len(name) > 0 && name[0] == '_' || name[0] == '.'
		}),
		scanner.Markdown(
			"content",
			func(path string, Metadata map[string]any) bool {
				return Metadata["draft"] != true
			},
			goldmark.WithExtensions(
				extension.GFM,
				extension.Footnote,
				highlighting.NewHighlighting(
					highlighting.WithStyle("monokai"),
					highlighting.WithFormatOptions(
						chromahtml.WithLineNumbers(true),
					),
				),
			),
			goldmark.WithRendererOptions(
			//	html.WithUnsafe(),
			),
		),
		scanner.Redirect(map[string]string{
			// legacy URLs
			"2016/02/01/curse-of-the-doctype": "/curse-of-the-doctype/",
			"2016/09/28/a-rreally-bad-idea":   "/a-rreally-bad-idea/",
			"2025/07/08/go-empty-struct":      "/go-empty-struct/",

			// published drafts
			"drafts/ggman": "/ggman/",
		}),
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

	Output: output.Native("public", true),
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

	// running with DEBUG=1 starts a server
	if os.Getenv("WATCH") != "" {
		var server http.Server
		server.Addr = "localhost:8080"

		g.Output, server.Handler = output.Server()

		done := make(chan error, 1)
		go func() {
			logger.Info("server listening", "addr", server.Addr)
			done <- server.ListenAndServe()
		}()

		if err := g.Watch(ctx, logger); err != nil {
			exitCode = 1
		}

		if exitCode != 0 {
			return
		}

		go func() {
			<-ctx.Done()
			server.Close()
		}()

		err := <-done
		logger.Info("server closed", "err", err)

		return
	}

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
