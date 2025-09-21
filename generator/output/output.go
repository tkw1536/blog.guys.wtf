// Package output provides Output.
//
//spellchecker:words generator
package output

//spellchecker:words context errors html template slog path filepath slices strings sync atomic
import (
	"context"
	"log/slog"

	"go.tkw01536.de/blog/generator/file"
)

// Output writes files into an abitrary location.
// Use [Native] for a default function.
type Output interface {
	// Write writes the given file into the output.
	// Write may be called concurrently.
	Write(ctx context.Context, logger *slog.Logger, file file.File) error

	// Reset is invoked right before the first file is written.
	// It may be used to initialize the output, or to reset it when a new generation occurs in Watch mode.
	Reset() error
}
