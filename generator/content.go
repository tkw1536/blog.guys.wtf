package generator

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log/slog"
)

// ContentFile is a file to be rendered as a content template.
type ContentFile struct {
	File
	Metadata map[string]any // Metadata contained in the file, if any.
}

func (cf *ContentFile) Body() template.HTML {
	return template.HTML(cf.File.Contents)
}

// IndexTemplate is a template to be used for generation.
type ContentTemplate struct {
	Template *template.Template // Template used for actual rendering, is passed [ContentTemplateContext].
	Globals  any                // Global Data to be passed.
}

func (ctc *ContentTemplate) Execute(w io.Writer, file ContentFile) error {
	return ctc.Template.Execute(w, &ContentTemplateContext{
		File:     file,
		Template: ctc,
	})
}

// ContentTemplateContext is passed to a [ContentTemplate].
type ContentTemplateContext struct {
	File     ContentFile // File to be rendered.
	Template *ContentTemplate
}

// Render the content pages.
func (generator *Generator) renderContents(
	ctx context.Context,
	logger *slog.Logger,

	output chan<- File,
	input <-chan ContentFile,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var out bytes.Buffer
	for file := range input {
		logger.Info("generating content file", slog.String("path", file.Path))

		out.Reset()
		if err := generator.ContentTemplate.Execute(&out, file); err != nil {
			return fmt.Errorf("failed to render content %q: %w", file.Path, err)
		}

		output <- File{
			Path:     file.Path,
			Contents: out.Bytes(),
		}

	}
	return nil
}
