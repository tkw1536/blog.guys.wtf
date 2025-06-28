package generator

import (
	"context"
	"io/fs"
	"log/slog"
)

// Scanner is a function that scans files.
// See [NewMarkdownScanner] and [NewStaticScanner].
type Scanner func(ctx context.Context, logger *slog.Logger, files chan<- ScannedFile) error

// ScannedFile is a file returned by a scanner.
type ScannedFile struct {
	File

	Metadata any  // Metadata contained in the file, if any.
	Indexed  bool // Should this file be indexed?

	Raw bool // if false, don't pass this file through the content template afterwards.
}

// NewStaticScanner does stuff
func NewStaticScanner(fsys fs.FS, excludes []string) Scanner {
	panic("not yet implemented")
}
