package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"regexp"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	mHTML "github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	gmHtml "github.com/yuin/goldmark/renderer/html"
	"golang.org/x/net/html"
)

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

func NewTemplate(src string, name string, common map[string]any) (*Template, error) {
	template, err := template.New(name).Funcs(templateFuncs).Parse(src)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &Template{
		tpl: template,
		markdown: goldmark.New(
			goldmark.WithRendererOptions(
				gmHtml.WithUnsafe(),
			),
			goldmark.WithExtensions(
				meta.Meta,
			),
		),
		common: common,
	}, nil

}

func MustTemplate(src, name string, common map[string]any) *Template {
	tpl, err := NewTemplate(src, name, common)
	if err != nil {
		panic(err)
	}
	return tpl
}

type Template struct {
	tpl      *template.Template
	markdown goldmark.Markdown
	common   map[string]any
}

type TplContext struct {
	Content template.HTML
	Meta    map[string]any
	Common  map[string]any
}

var m *minify.M

func init() {
	m = minify.New()
	m.AddFunc("text/html", mHTML.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
}

// Renders the template, returning content and metadata
func (tpl *Template) Render(content []byte, metaIn map[string]any) ([]byte, map[string]any, error) {
	context := parser.NewContext()

	// no metadata => render in markdown
	if metaIn == nil {
		var mdBuffer bytes.Buffer
		if err := tpl.markdown.Convert(content, &mdBuffer, parser.WithContext(context)); err != nil {
			return nil, nil, fmt.Errorf("failed to convert markdown: %w", err)
		}

		var outBuffer bytes.Buffer
		if err := addRelToLinks(&outBuffer, &mdBuffer); err != nil {
			return nil, nil, fmt.Errorf("failed to make links open in new tab: %w", err)
		}

		metaIn = meta.Get(context)
		content = outBuffer.Bytes()
	}

	var htmlBuffer bytes.Buffer
	if err := tpl.tpl.Execute(&htmlBuffer, TplContext{
		Content: template.HTML(content),
		Meta:    metaIn,
		Common:  tpl.common,
	}); err != nil {
		return nil, nil, fmt.Errorf("failed to render template: %w", err)
	}

	var minifyBuffer bytes.Buffer
	if err := m.Minify("text/html", &minifyBuffer, &htmlBuffer); err != nil {
		return nil, nil, fmt.Errorf("failed to minify: %w", err)
	}

	return minifyBuffer.Bytes(), metaIn, nil
}

// addRelToLinks adds rel="noopener noreferrer" to all links in the given HTML string.
func addRelToLinks(dst io.Writer, src io.Reader) error {

	process := func(n *html.Node) {
		if n.Type != html.ElementNode || n.Data != "a" {
			return
		}

		var (
			hrefId   = -1
			relId    = -1
			targetId = -1
		)
		for i, attr := range n.Attr {
			switch attr.Key {
			case "href":
				hrefId = i
			case "rel":
				relId = i
			case "target":
				targetId = i
			}
			if attr.Key == "href" {
				hrefId = i
			}
		}

		if hrefId == -1 {
			return
		}

		if relId >= 0 {
			n.Attr[relId].Val += " noopener noreferrer"
		} else {
			n.Attr = append(n.Attr, html.Attribute{Key: "rel", Val: "noopener noreferrer"})
		}

		if targetId >= 0 {
			n.Attr[targetId].Val += "_blank"
		} else {
			n.Attr = append(n.Attr, html.Attribute{Key: "target", Val: "blank"})
		}
	}

	// Walk through the HTML document and modify links
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		process(n)
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	root, err := html.Parse(src)
	if err != nil {
		return fmt.Errorf("failed to parse html: %w", err)
	}

	walk(root)

	if err := html.Render(dst, root); err != nil {
		return fmt.Errorf("failed to render modified document: %w", err)
	}
	return nil
}
