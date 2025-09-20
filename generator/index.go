//spellchecker:words generator
package generator

//spellchecker:words bytes context html template slog slices strings
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

	CompareFunc IndexComparisonFunc

	Template *template.Template // Template to use for rendering.
	Globals  map[string]any     // Global Metadata
	Metadata map[string]any     // Metadata to return from the template.
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

// IndexEntry represents a single indexed page.
type IndexEntry struct {
	Path     string         // Path the file will be outputted in
	Metadata map[string]any // Metadata contained in the file, if any
}

// Link returns a nice link to this page.
func (index IndexEntry) Link() string {
	return File{Path: index.Path}.Link()
}

// An IndexComparisonFunc is passed to [slices.SortFunc] to compare to indexes.
// It is invoked in order to establish a total order on index entries.
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

	for _, tpl := range generator.Indexes {
		logger.Info("sorting index", slog.Int("entryCount", len(entries)))
		slices.SortFunc(entries, tpl.CompareFunc.f())

		logger.Info("generating index content", slog.String("path", tpl.Path), slog.Int("entryCount", len(entries)))

		var out bytes.Buffer
		if err := tpl.Execute(&out, entries); err != nil {
			return fmt.Errorf("failed to render index contents %q: %w", tpl.Path, err)
		}

		output <- ScannedFile{
			FileWithMetadata: FileWithMetadata{
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
