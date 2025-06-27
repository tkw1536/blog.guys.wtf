package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// File represents a single file to be placed in the output directory.
type File struct {
	Path     string
	Contents []byte
}

// outputFiles writes files into the dest directory
func outputFiles(ctx context.Context, logger *log.Logger, dest string, files <-chan File) error {
	if err := ensureEmptyDir(dest); err != nil {
		return fmt.Errorf("failed to ensure empty directory %q: %w", dest, err)
	}
	for {
		select {
		case file, ok := <-files:
			if !ok { // no more files
				return nil
			}

			dest := filepath.Join(dest, file.Path)
			dir := filepath.Dir(dest)

			logger.Printf("creating directory %q", dir)
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory for %q: %w", file.Path, err)
			}

			logger.Printf("writing file %q", dest)
			if err := os.WriteFile(dest, file.Contents, os.ModePerm); err != nil {
				return fmt.Errorf("failed to write file %q: %w", dest, err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// ensureEmptyDir ensures that dest exists and is an empty directory.
func ensureEmptyDir(dest string) error {
	info, err := os.Stat(dest)

	// doesn't exist => create it
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dest, os.ModePerm); err != nil {
			return fmt.Errorf("failed to make directory: %w", err)
		}
		return nil
	}

	if !info.IsDir() {
		return fmt.Errorf("not a directory: %w", err)
	}

	// delete all content of the directory!
	d, err := os.Open(dest)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return fmt.Errorf("failed to list contents for deletion: %w", err)
	}
	for _, name := range names {
		content := filepath.Join(dest, name)
		err = os.RemoveAll(content)
		if err != nil {
			return fmt.Errorf("failed to remove %q: %w", content, err)
		}
	}
	return nil
}
