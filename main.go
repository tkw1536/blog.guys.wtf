package main

import (
	"cmp"
	"context"
	"generator/generator"
	"html/template"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	_ "embed"
)

var commonData = map[string]any{
	"BlogTitle": "High on Code!",
}

//go:embed "templates/index.html"
var indexHTML string
var htmlTemplate = generator.MustTemplate(indexHTML, "index.html", commonData)

//go:embed templates/list.html
var listHTML string
var listTemplate = generator.MustPlainTemplate(listHTML, "list.html")

var g = generator.Generator{
	Contents: os.DirFS("content"),
	Static:   os.DirFS("static"),

	Indexes: map[string]*template.Template{
		"index.html": listTemplate,
	},
	IndexCompareFunc: func(left, right generator.Indexed) int {
		lMeta, _ := left.Meta.(map[string]any)
		lDate, _ := lMeta["date"].(string)

		rMeta, _ := right.Meta.(map[string]any)
		rDate, _ := rMeta["date"].(string)

		return cmp.Or(
			strings.Compare(rDate, lDate),
			strings.Compare(left.Path, right.Path),
		)
	},

	ContentTemplate: htmlTemplate,

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
