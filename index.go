package main

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"strings"

	_ "embed"
)

// Indexed is an indexed blog page.
type Indexed struct {
	Path string
	Meta map[string]any
}

// Link returns a nice link to this page.
func (index Indexed) Link() string {
	noIndex := strings.TrimSuffix(index.Path, "/index.html")
	return strings.TrimSuffix(noIndex, "/") + "/"
}

//go:embed templates/list.html
var listHTML string
var listTemplate = template.Must(template.New("list.html").Funcs(templateFuncs).Parse(listHTML))

// Render the index page
func RenderIndex(ctx context.Context, logger *log.Logger, indexed []Indexed, template *Template) File {
	var out bytes.Buffer
	if err := listTemplate.Execute(&out, indexed); err != nil {
		log.Fatalf("failed to render index: %v", err)
	}

	contents, _, err := template.Render(out.Bytes(), map[string]any{})
	if err != nil {
		log.Fatalf("failed to render index: %v", err)
	}

	return File{
		Path:     "index.html",
		Contents: contents,
	}
}
