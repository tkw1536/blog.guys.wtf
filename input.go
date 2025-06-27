package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	errFileExcluded = errors.New("file excluded")
	errNotMarkdown  = errors.New("not a markdown file")
)

func scanStatic(ctx context.Context, logger *log.Logger, output chan<- File, input string) error {
	return scanAndProcess(
		ctx, logger, output,
		func(info fs.FileInfo) error {
			name := info.Name()
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
		input,
	)
}

func scanMarkdown(ctx context.Context, logger *log.Logger, template *Template, output chan<- File, index chan<- Indexed, input string) error {
	return scanAndProcess(
		ctx, logger, output,
		func(info fs.FileInfo) error {
			name := info.Name()
			if !strings.HasSuffix(name, ".md") {
				return errFileExcluded
			}
			return nil
		},
		func(path string, fi fs.FileInfo, b []byte) (File, error) {
			contents, meta, err := template.Render(b, nil)
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
		input,
	)
}

func scanAndProcess(ctx context.Context, logger *log.Logger, output chan<- File, include func(info fs.FileInfo) error, process func(string, fs.FileInfo, []byte) (File, error), input string) error {
	if err := filepath.Walk(input, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// don't do anything with directories
		if info.IsDir() {
			return nil
		}

		// check if we should include the file
		if err := include(info); err != nil {
			logger.Printf("excluding %q: %v", path, err)
			return nil
		}

		rel, err := filepath.Rel(input, path)
		if err != nil {
			return fmt.Errorf("failed to find relative path for %q: %w", input, err)
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read contents of %q: %w", input, err)
		}

		logger.Printf("processing %q", path)
		processed, err := process(rel, info, contents)
		if err != nil {
			return fmt.Errorf("failed to process contents of %q: %w", input, err)
		}
		output <- processed

		return nil
	}); err != nil {
		return fmt.Errorf("failed to scan files: %w", err)
	}

	return nil
}
