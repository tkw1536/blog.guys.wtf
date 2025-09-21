// Package scanner holds scanners that provide input to the generator.
//
//spellchecker:words generator
package scanner

//spellchecker:words context errors slog
import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"

	"go.tkw01536.de/blog/generator/file"
)

// Scanner scans abitrary sources for inputs to the generator.
// To create a new scanner, use [Markdown], [Static] or using [Redirect].
type Scanner interface {
	// Scan scans the underlying source and sends inputs to the given channel.
	//
	// If a file is to be excluded, it may return an error that wraps [ErrExcluded].
	// If any other error occurs, it may be treated as fatal by the caller.
	Scan(ctx context.Context, logger *slog.Logger, files chan<- file.ScannedFile) error

	// Paths returns a list of file system paths this input scanner depends on.
	// These are used to watch for changes, and automatically trigger a rebuild.
	Paths() []string
}

// ErrExcluded is used by [newFSScanner] to indicate that a file is to be skipped.
var ErrExcluded = errors.New("file excluded")

// fsScanner is a scanner operating on a filesystem.
type fsScanner struct {
	// open opens the given filesystem.
	open func() (fs.FS, error)
	// process processes a single file from the filesystem into a file.
	process func(path string, d fs.DirEntry, contents []byte) (file.ScannedFile, error)
	paths   []string
}

func (scanner *fsScanner) Scan(ctx context.Context, logger *slog.Logger, files chan<- file.ScannedFile) error {
	fsys, err := scanner.open()
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

		file, err := scanner.process(path, d, contents)
		if errors.Is(err, ErrExcluded) {
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

func (scanner *fsScanner) Paths() []string {
	return scanner.paths
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
