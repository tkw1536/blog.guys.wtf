//spellchecker:words generator
package generator

//spellchecker:words context errors html template slog path filepath slices strings sync atomic
import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"sync/atomic"
)

// FileWriter is a function that writes a file to output.
// Use [NewNativeFileWriter] for a default function.
type FileWriter interface {
	Reset() error
	Write(ctx context.Context, logger *slog.Logger, file File) error
}

// NewNativeFileWriter creates a new [FileWriter] that writes to the given root as its' output directory.
//
// Files outside the given directory are not tracked.
// If cleanFirst is set to true, it cleans all files when first invoked.
func NewNativeFileWriter(path string, cleanFirst bool) FileWriter {
	return &nativeFileWriter{
		path:       path,
		cleanFirst: cleanFirst,
	}
}

type nativeFileWriter struct {
	loaded atomic.Bool

	m    sync.Mutex
	root *os.Root

	path       string
	cleanFirst bool
}

func (nfw *nativeFileWriter) openRoot(logger *slog.Logger) (*os.Root, error) {
	if !nfw.loaded.Load() {
		return nfw.openRootSlow(logger)
	}

	return nfw.root, nil
}

func (nfw *nativeFileWriter) openRootSlow(logger *slog.Logger) (*os.Root, error) {
	nfw.m.Lock()
	defer nfw.m.Unlock()

	// already loaded
	if nfw.root != nil {
		return nfw.root, nil
	}

	root, err := openRoot(logger, nfw.path, nfw.cleanFirst)
	if err != nil {
		return nil, fmt.Errorf("openRoot: %w", err)
	}

	nfw.root = root
	nfw.loaded.Store(true)
	return nfw.root, nil
}

func (nfw *nativeFileWriter) Reset() error {
	return nil
}

func (nfw *nativeFileWriter) Write(ctx context.Context, logger *slog.Logger, file File) (e error) {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context closed: %w", err)
	}

	// open the root directory
	root, err := nfw.openRoot(logger)
	if err != nil {
		return fmt.Errorf("failed to open root directory: %w", err)
	}

	path := file.Path
	parent := filepath.Dir(path)

	if err := mkdirAll(root, parent, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	logger.Info("writing file", slog.String("path", path), slog.Int("size", len(file.Contents)))
	handle, err := root.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		errClose := handle.Close()
		if errClose == nil {
			return
		}
		errClose = fmt.Errorf("failed to close file handle: %w", errClose)

		if e == nil {
			e = errClose
		} else {
			e = errors.Join(e, errClose)
		}
	}()

	if _, err := handle.Write(file.Contents); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// openRoot is like [os.OpenRoot], except that it possibly creates path if it doesn't exist, and optionally removes any existing files in it.
func openRoot(logger *slog.Logger, path string, clean bool) (*os.Root, error) {
	root, err := os.OpenRoot(path)

	// if it doesn't exist => create it!
	if os.IsNotExist(err) {
		logger.Info("creating directory", slog.String("path", path))
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, fmt.Errorf("os.MkdirAll: %w", err)
		}

		root, err = os.OpenRoot(path)
	}
	if err != nil {
		return nil, fmt.Errorf("os.OpenRoot: %w", err)
	}

	if !clean {
		return root, nil
	}

	logger.Info("removing contents of directory", slog.String("path", path))
	if err := removeAllContents(root, "."); err != nil {
		return nil, fmt.Errorf("removeAllContents: %w", err)
	}

	return root, nil
}

func mkdirAll(root *os.Root, path string, perm os.FileMode) error {
	// parent directories yet to create
	var paths []string

	oldPath := path + "/" // first iteration of the loop must always happen!
	for path != oldPath {

		// check existence so we can bail out early.
		_, err := root.Stat(path)
		if os.IsExist(err) {
			break
		}
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("root.Stat(%q): %w", path, err)
		}

		// append the path
		paths = append(paths, path)

		// go to the next path
		oldPath = path
		path = filepath.Dir(path)
	}

	for _, name := range slices.Backward(paths) {
		err := root.Mkdir(name, perm)
		if os.IsExist(err) {
			err = nil // directory already exists
		}
		if err != nil {
			return fmt.Errorf("root.Mkdir(%q): %w", name, err)
		}
	}
	return nil
}

// removeAll is like [os.Remove] but works on an *os.Root.
func removeAll(root *os.Root, path string) error {
	f, err := root.Stat(path)
	if err != nil {
		return fmt.Errorf("root.Stat: %w", err)
	}

	if f.IsDir() {
		if err := removeAllContents(root, path); err != nil {
			return err
		}
	}

	if err := root.Remove(path); err != nil {
		return fmt.Errorf("root.Remove: %w", err)
	}
	return nil
}

// removeAllContents removes all contents of the given directory.
// The directory itself is not removed.
func removeAllContents(root *os.Root, path string) error {
	dir, err := root.Open(path)
	if err != nil {
		return fmt.Errorf("root.Open(): %w", err)
	}
	defer dir.Close()

	entries, err := dir.ReadDir(-1)
	if err != nil {
		return fmt.Errorf("dir.ReadDir: %w", err)
	}

	for _, entry := range entries {
		path := filepath.Join(path, entry.Name())
		if err := removeAll(root, path); err != nil {
			return err
		}
	}

	return nil
}
