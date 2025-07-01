package generator

import (
	"io/fs"
	"strings"
)

// NewStaticScanner adds a scanner that copies files from path.
// The found files are not added to the index.
//
// If files start with any of prefixes, it is ignored.
// Scanner internally uses [os.Root], and ensures that no files outside the given directory are caught.
func NewStaticScanner(path string, prefixes []string) Scanner {
	return newFSScanner(
		openRootFS(path),
		func(path string, d fs.DirEntry, contents []byte) (ScannedFile, error) {
			// check if the file is excluded
			name := d.Name()
			for _, exclude := range prefixes {
				if strings.HasPrefix(name, exclude) {
					return ScannedFile{}, errExcluded
				}
			}

			// then copy it as-is
			return ScannedFile{
				FileWithMetadata: FileWithMetadata{
					File: File{
						Path:     path,
						Contents: contents,
					},
					Metadata: nil,
				},
				Indexed: false,
				Raw:     true,
			}, nil
		},
	)
}
