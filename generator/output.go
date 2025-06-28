package generator

import (
	"context"
	"fmt"
	"log/slog"
)

// outputFiles writes files to the given output
func (generate *Generator) outputFiles(ctx context.Context, logger *slog.Logger, files <-chan File) error {
	for {
		select {
		case file, ok := <-files:
			if !ok { // no more files
				return nil
			}

			if err := generate.Output(ctx, logger, file); err != nil {
				return fmt.Errorf("failed to output file: %w", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
