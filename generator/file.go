package generator

import (
	"html/template"
	"strings"
)

// File describes a single file.
type File struct {
	// Path is the relative path relative to some abitrary root directory to this file.
	//
	// Paths may start with ".." indicating behavior outside the root output directory.
	// Consumers of a File should implement appropriate protections.
	Path string

	Contents []byte
}

// Body returns the contents of this file as unsafe html.
// Intended to be used in templates.
func (cf *File) Body() template.HTML {
	return template.HTML(cf.Contents)
}

// Link returns a nice link to this file.
// Links always start with "/", and only end in slash in case of a directory.
func (file File) Link() string {
	if file.Path == "index.html" {
		return "/"
	}

	cleanPath := strings.Trim(file.Path, "/")

	if strings.HasSuffix(cleanPath, "/index.html") {
		cleanPath = cleanPath[:len(cleanPath)-len("index.html")]
	}

	return "/" + cleanPath
}

// FileWithMetadata represents a file along with associated metadata.
type FileWithMetadata struct {
	File
	Metadata map[string]any // Metadata contained in the file, if any.
}

// ScannedFile represents a file returned by a scanner.
type ScannedFile struct {
	FileWithMetadata

	Indexed bool // Should this file be indexed?
	Raw     bool // if false, don't pass this file through the content template afterwards.
}
