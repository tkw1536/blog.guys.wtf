package generator

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

// Indexed is an indexed blog page.
type Indexed struct {
	Path string
	Meta any
}

// Link returns a nice link to this page.
func (index Indexed) Link() string {
	noIndex := strings.TrimSuffix(index.Path, "/index.html")
	return strings.TrimSuffix(noIndex, "/") + "/"
}

// An IndexComparisonFunc is passed to [slices.SortFunc] to compare to indexes.
type IndexComparisonFunc func(left, right Indexed) int

// f returns the function used to sort
func (ifc IndexComparisonFunc) f() func(left, right Indexed) int {
	if ifc == nil {
		return func(left, right Indexed) int {
			return strings.Compare(left.Path, right.Path)
		}
	}
	return ifc
}

// Render the index pages.
func (generator *Generator) RenderIndexes(
	ctx context.Context,
	logger *slog.Logger,

	indexed []Indexed,
	output chan<- File,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	slices.SortFunc(indexed, generator.IndexCompareFunc.f())

	for path, tpl := range generator.Indexes {
		logger.Info("generating index page contents", slog.Int("content", len(indexed)), slog.Int("size", len(indexed)))
		var out bytes.Buffer
		if err := tpl.Execute(&out, indexed); err != nil {
			return fmt.Errorf("failed to render index contents %q: %w", path, err)
		}

		logger.Info("rendering index page", slog.String("path", path))
		contents, _, err := generator.ContentTemplate.Render(out.Bytes(), map[string]any{})
		if err != nil {
			return fmt.Errorf("failed to render index %q: %w", path, err)
		}

		output <- File{
			Path:     path,
			Contents: contents,
		}
	}
	return nil
}
