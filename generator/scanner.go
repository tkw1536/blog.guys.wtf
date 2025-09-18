//spellchecker:words generator
package generator

//spellchecker:words context errors slog
import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
)

// Scanner is a function that scans for inputs.
// See [NewMarkdownScanner] and [NewStaticScanner].
type Scanner struct {
	scan  scannerFunc
	paths []string
}

type scannerFunc = func(ctx context.Context, logger *slog.Logger, files chan<- ScannedFile) error

// ScannedFile is a file returned by a scanner.
type ScannedFile struct {
	FileWithMetadata

	Indexed bool // Should this file be indexed?
	Raw     bool // if false, don't pass this file through the content template afterwards.
}

// errExcluded is used by [newFSScanner] to indicate that a file is to be skipped.
var errExcluded = errors.New("file excluded")

func newFSScanner(open func() (fs.FS, error), process func(path string, d fs.DirEntry, contents []byte) (ScannedFile, error)) scannerFunc {
	return func(ctx context.Context, logger *slog.Logger, files chan<- ScannedFile) error {
		fsys, err := open()
		if err != nil {
			return fmt.Errorf("failed to open filesystem: %w", err)
		}

		if err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if err := ctx.Err(); err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			contents, err := fs.ReadFile(fsys, path)
			if err != nil {
				return fmt.Errorf("failed to read file %q: %w", path, err)
			}

			file, err := process(path, d, contents)
			if errors.Is(err, errExcluded) {
				logger.Info("skipping file %q", slog.String("path", path))
				return nil
			}

			if err != nil {
				return fmt.Errorf("failed to process file %q: %w", path, err)
			}

			logger.Info("scanned file", slog.String("path", path))
			files <- file
			return nil
		}); err != nil {
			return fmt.Errorf("WalkDir failed: %w", err)
		}
		return nil
	}
}

// openRootFS is like [os.OpenRoot] followed by [os.Root.FS].
func openRootFS(path string) func() (fs.FS, error) {
	return func() (fs.FS, error) {
		root, err := os.OpenRoot(path)
		if err != nil {
			return nil, fmt.Errorf("os.OpenRoot: %w", err)
		}
		return root.FS(), nil
	}
}
