package generator

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

// PostProcessor post processes a file before outputting
type PostProcessor func(in File) (out File, err error)

func (generator *Generator) postProcess(
	ctx context.Context,
	logger *slog.Logger,
	file File,
) (File, error) {
	logger.Info("post processing file", slog.String("path", file.Path))
	for _, processor := range generator.PostProcessors {
		var err error
		file, err = processor(file)
		if err != nil {
			return File{}, fmt.Errorf("failed to post process: %w", err)
		}
	}

	return file, nil
}

var m *minify.M

func init() {
	m = minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
}

func MinifyPostProcessor(in File) (out File, err error) {
	ext := strings.ToLower(filepath.Ext(in.Path))

	var mediaType string
	switch ext {
	case ".css":
		mediaType = "text/css"
	case ".html", ".htm":
		mediaType = "text/html"
	case ".svg":
		mediaType = "image/svg+xml"
	case ".js", ".mjs":
		mediaType = "text/javascript"
	case ".json":
		mediaType = "text/json"
	case "xml":
		mediaType = "text/xml"
	default:
		return in, nil
	}

	var buffer bytes.Buffer
	if err := m.Minify(mediaType, &buffer, bytes.NewReader(in.Contents)); err != nil {
		return File{}, fmt.Errorf("failed to minify %q: %w", in.Path, err)
	}
	return File{
		Path:     in.Path,
		Contents: buffer.Bytes(),
	}, nil
}

var _ PostProcessor = MinifyPostProcessor
