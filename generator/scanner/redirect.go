//spellchecker:words generator
package scanner

//spellchecker:words bytes context html template slog strings
import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"strings"

	"go.tkw01536.de/blog/generator/file"
)

var redirectTemplate = template.Must(template.New("").Parse(`<!DOCTYPE html><title>{{ . }}</title><meta http-equiv="refresh" content = "0;url={{.}}" />`))

// Redirect adds static html files that redirect from source to target.
func Redirect(sourceToTarget map[string]string) Scanner {
	return redirectScanner(sourceToTarget)
}

type redirectScanner map[string]string

func (scanner redirectScanner) Scan(ctx context.Context, logger *slog.Logger, files chan<- file.ScannedFile) error {
	for source, target := range scanner {
		var buffer bytes.Buffer
		if err := redirectTemplate.Execute(&buffer, target); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		path := strings.Trim(source, "/") + "/index.html"

		file := file.ScannedFile{
			FileWithMetadata: file.FileWithMetadata{
				File: file.File{Path: path, Contents: buffer.Bytes()},
			},
			Raw: true,
		}

		select {
		case files <- file:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (scanner redirectScanner) Paths() []string {
	return nil
}
