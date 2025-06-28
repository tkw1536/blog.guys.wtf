package generator

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
)

var (
	errFileExcluded = errors.New("file excluded")
	errNotMarkdown  = errors.New("not a markdown file")
)

func (generator *Generator) scanStatic(
	ctx context.Context, logger *slog.Logger,
	output chan<- File,
) error {
	return scanAndProcess(
		ctx, logger, output,

		generator.Static,
		func(path string, info fs.FileInfo) error {
			name := info.Name()
			// TODO: Make this generic
			if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
				return errFileExcluded
			}
			return nil
		},
		func(path string, fi fs.FileInfo, b []byte) (File, error) {
			return File{
				Path:     path,
				Contents: b,
			}, nil
		},
	)
}

func (generator *Generator) scanMarkdown(
	ctx context.Context, logger *slog.Logger,
	output chan<- File,
	index chan<- Indexed,
) error {
	return scanAndProcess(
		ctx, logger, output,

		generator.Contents,
		func(path string, info fs.FileInfo) error {
			name := info.Name()
			if !strings.HasSuffix(name, ".md") {
				return errFileExcluded
			}
			return nil
		},
		func(path string, fi fs.FileInfo, b []byte) (File, error) {
			contents, meta, err := generator.ContentTemplate.Render(b, nil)
			if err != nil {
				return File{}, fmt.Errorf("failed to render: %w", err)
			}

			// by default, make the destination file '[slug]/index.html'
			filename := filepath.Join(path[:len(path)-len(".md")], "index.html")

			// if we have _[something].md directly output that as [something].html
			if name := fi.Name(); strings.HasPrefix(name, "_") {
				nameWithHTML := name[:len(name)-len(".md")] + ".html"
				nameWithHTML = nameWithHTML[1:]
				filename = filepath.Join(filepath.Dir(path), nameWithHTML)
			}

			index <- Indexed{
				Path: filename,
				Meta: meta,
			}

			return File{
				Path:     filename,
				Contents: contents,
			}, nil
		},
	)
}

func scanAndProcess(
	ctx context.Context, logger *slog.Logger,
	output chan<- File,

	fsys fs.FS,
	include func(path string, info fs.FileInfo) error,
	process func(string, fs.FileInfo, []byte) (File, error),
) error {

	if err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}

		// don't do anything with dirs!
		if d.IsDir() {
			return nil
		}

		// get file info
		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("failed to get info on file: %w", err)
		}

		// if it's not included => don't do anything
		if err := include(path, info); err != nil {
			logger.Info("skipping file", slog.String("path", path))
			return nil
		}

		// read the file's contents
		contents, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("failed to read contents of %q: %w", path, err)
		}

		// and do the processing
		logger.Info("processing", slog.String("path", path))
		processed, err := process(path, info, contents)
		if err != nil {
			return fmt.Errorf("failed to process contents of %q: %w", path, err)
		}
		output <- processed

		return nil

	}); err != nil {
		return fmt.Errorf("failed to process files: %w", err)
	}

	return nil
}
