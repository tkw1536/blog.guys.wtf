package generator

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"slices"
	"strings"
)

// IndexTemplate is an index file to be generated.
type IndexTemplate struct {
	Path string // Path of the file to create.
	Raw  bool   // Should the path assumed to be raw?

	Template *template.Template // Template to use for rendering.
	Globals  any                // Global Metadata
	Metadata any                // Metadata to return from the template.
}

// Execute executes this index template.
func (tpl *IndexTemplate) Execute(w io.Writer, entries []IndexEntry) error {
	if err := tpl.Template.Execute(w, &IndexTemplateContext{
		Entries:  entries,
		Template: tpl,
	}); err != nil {
		return fmt.Errorf("failed to execute index template: %w", err)
	}
	return nil
}

// IndexTemplateContext is passed to an index template.
type IndexTemplateContext struct {
	Entries  []IndexEntry
	Template *IndexTemplate
}

// IndexEntry is an indexed blog page.
type IndexEntry struct {
	Path     string // Path the file will be outputted in
	Metadata any    // Metadata contained in the file, if any
}

// Link returns a nice link to this page.
func (index IndexEntry) Link() string {
	noIndex := strings.TrimSuffix(index.Path, "/index.html")
	return strings.TrimSuffix(noIndex, "/") + "/"
}

// An IndexComparisonFunc is passed to [slices.SortFunc] to compare to indexes.
type IndexComparisonFunc func(left, right IndexEntry) int

// f returns the function used to sort
func (ifc IndexComparisonFunc) f() func(left, right IndexEntry) int {
	if ifc == nil {
		return func(left, right IndexEntry) int {
			return strings.Compare(left.Path, right.Path)
		}
	}
	return ifc
}

// Render the index pages.
func (generator *Generator) renderIndexes(
	ctx context.Context,
	logger *slog.Logger,

	entries []IndexEntry,
	output chan<- ScannedFile,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	slices.SortFunc(entries, generator.IndexCompareFunc.f())

	var out bytes.Buffer
	for _, tpl := range generator.Indexes {
		logger.Info("generating index content", slog.String("path", tpl.Path), slog.Int("entryCount", len(entries)))

		out.Reset()
		if err := tpl.Execute(&out, entries); err != nil {
			return fmt.Errorf("failed to render index contents %q: %w", tpl.Path, err)
		}

		output <- ScannedFile{
			ContentFile: ContentFile{
				File: File{
					Path:     tpl.Path,
					Contents: out.Bytes(),
				},
				Metadata: tpl.Metadata,
			},
			Indexed: false,
			Raw:     tpl.Raw,
		}
	}
	return nil
}
