//spellchecker:words generator
package generator

//spellchecker:words bytes context html template slog strings
import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"strings"
)

var redirectTemplate = template.Must(template.New("").Parse(`<!DOCTYPE html><title>{{ . }}</title><meta http-equiv="refresh" content = "0;url={{.}}" />`))

// NewRedirectScanner adds static html files that redirect from source to target
func NewRedirectScanner(sourceToTarget map[string]string) *Scanner {
	return &Scanner{
		scan: func(ctx context.Context, logger *slog.Logger, files chan<- ScannedFile) error {

			for source, target := range sourceToTarget {
				var buffer bytes.Buffer
				if err := redirectTemplate.Execute(&buffer, target); err != nil {
					return fmt.Errorf("failed to execute template: %w", err)
				}

				path := strings.Trim(source, "/") + "/index.html"

				file := ScannedFile{
					FileWithMetadata: FileWithMetadata{
						File: File{Path: path, Contents: buffer.Bytes()},
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
		},
		paths: nil,
	}
}
