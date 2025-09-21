//spellchecker:words generator
package scanner

//spellchecker:words strings
import (
	"io/fs"

	"go.tkw01536.de/blog/generator/file"
)

// Static adds a scanner that copies files from path into the output.
// The found files are not added to the index.
//
// If exclude is not nil and exclude(file.Name()) returns true, it is excluded.
//
// Scanner internally uses [os.Root], and ensures that no files outside the given directory are copied.
func Static(path string, exclude func(name string) bool) Scanner {
	return &fsScanner{
		open: openRootFS(path),
		process: func(path string, d fs.DirEntry, contents []byte) (file.ScannedFile, error) {
			// check if the file is to be excluded
			if exclude != nil && exclude(d.Name()) {
				return file.ScannedFile{}, ErrExcluded
			}

			// then copy it as-is
			return file.ScannedFile{
				FileWithMetadata: file.FileWithMetadata{
					File: file.File{
						Path:     path,
						Contents: contents,
					},
					Metadata: nil,
				},
				Indexed: false,
				Raw:     true,
			}, nil
		},
		paths: []string{path},
	}
}
