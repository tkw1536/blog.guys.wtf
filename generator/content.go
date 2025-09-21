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

	"go.tkw01536.de/blog/generator/file"
)

// IndexTemplate is a template to be used for generation.
type ContentTemplate struct {
	Template *template.Template // Template used for actual rendering, is passed [ContentTemplateContext].
	Globals  any                // Global Data to be passed.
}

func (ctc *ContentTemplate) Execute(w io.Writer, file file.FileWithMetadata) error {
	return ctc.Template.Execute(w, &ContentTemplateContext{
		File:     file,
		Template: ctc,
	})
}

// ContentTemplateContext is passed to a [ContentTemplate].
type ContentTemplateContext struct {
	File     file.FileWithMetadata // File to be rendered.
	Template *ContentTemplate
}

// renderFile renders a single [FileWithMetadata] through the [ContentTemplate]
func (generator *Generator) renderFile(ctx context.Context, logger *slog.Logger, f file.FileWithMetadata) (file.File, error) {
	logger.Info("generating content file", slog.String("path", f.Path))

	var out bytes.Buffer
	if err := generator.ContentTemplate.Execute(&out, f); err != nil {
		return file.File{}, fmt.Errorf("failed to render content %q: %w", f.Path, err)
	}

	return file.File{
		Path:     f.Path,
		Contents: out.Bytes(),
	}, nil
}
