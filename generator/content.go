//spellchecker:words generator
package generator

//spellchecker:words bytes context html template slog
import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log/slog"
)

// FileWithMetadata is a file with associated metadata.
type FileWithMetadata struct {
	File
	Metadata map[string]any // Metadata contained in the file, if any.
}

// IndexTemplate is a template to be used for generation.
type ContentTemplate struct {
	Template *template.Template // Template used for actual rendering, is passed [ContentTemplateContext].
	Globals  any                // Global Data to be passed.
}

func (ctc *ContentTemplate) Execute(w io.Writer, file FileWithMetadata) error {
	return ctc.Template.Execute(w, &ContentTemplateContext{
		File:     file,
		Template: ctc,
	})
}

// ContentTemplateContext is passed to a [ContentTemplate].
type ContentTemplateContext struct {
	File     FileWithMetadata // File to be rendered.
	Template *ContentTemplate
}

// renders a single content page
func (generator *Generator) renderContent(ctx context.Context, logger *slog.Logger, file FileWithMetadata) (File, error) {
	logger.Info("generating content file", slog.String("path", file.Path))

	var out bytes.Buffer
	if err := generator.ContentTemplate.Execute(&out, file); err != nil {
		return File{}, fmt.Errorf("failed to render content %q: %w", file.Path, err)
	}

	return File{
		Path:     file.Path,
		Contents: out.Bytes(),
	}, nil
}
