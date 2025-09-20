//spellchecker:words generator
package generator

//spellchecker:words bytes path filepath strings github yuin goldmark meta parser golang html
import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/net/html"
)

// NewMarkdownScanner adds a scanner that renders file with the "md" extension at path as markdown.
// files are added to the index if the index function returns true, or shouldIndex is nil.
//
// Internally uses [os.Root], and ensures that no files outside the given directory are caught.
func NewMarkdownScanner(path string, shouldIndex func(path string, Metadata map[string]any) bool, options ...goldmark.Option) Scanner {
	markdown := goldmark.New(
		append([]goldmark.Option{
			goldmark.WithExtensions(
				meta.Meta,
			),
		}, options...)...,
	)
	return &fsScanner{
		open: openRootFS(path),
		process: func(path string, d fs.DirEntry, contents []byte) (ScannedFile, error) {
			// check if the file is excluded
			name := d.Name()
			if !strings.HasSuffix(name, ".md") {
				return ScannedFile{}, errExcluded
			}

			context := parser.NewContext()

			// parse markdown
			var markdownResult bytes.Buffer
			if err := markdown.Convert(contents, &markdownResult, parser.WithContext(context)); err != nil {
				return ScannedFile{}, fmt.Errorf("failed to convert markdown: %w", err)
			}

			// addRel to external links
			var contentBuffer bytes.Buffer
			if err := addTargetAndRel(&contentBuffer, &markdownResult); err != nil {
				return ScannedFile{}, fmt.Errorf("failed to make links open in new tab: %w", err)
			}

			// check if we should index!
			metadata := meta.Get(context)
			doIndex := true
			if shouldIndex != nil {
				doIndex = shouldIndex(path, metadata)
			}

			// by default, make the destination file '[slug]/index.html'
			filename := filepath.Join(path[:len(path)-len(".md")], "index.html")

			// if we have _[something].md directly output that as [something].html
			if name := d.Name(); strings.HasPrefix(name, "_") {
				nameWithHTML := name[:len(name)-len(".md")] + ".html"
				nameWithHTML = nameWithHTML[1:]
				filename = filepath.Join(filepath.Dir(path), nameWithHTML)
			}

			// and then use
			return ScannedFile{
				FileWithMetadata: FileWithMetadata{
					File: File{
						Path:     filename,
						Contents: contentBuffer.Bytes(),
					},
					Metadata: metadata,
				},
				Indexed: doIndex,
				Raw:     false,
			}, nil
		},
		paths: []string{path},
	}
}

// addTargetAndRel adds target="_blank" rel="noopener noreferrer" to all links in the given HTML, unless they start with '#'
func addTargetAndRel(dst io.Writer, src io.Reader) error {
	// updateLink updates a token representing an '<a' starting element.
	updateLink := func(token *html.Token) {
		var (
			hrefId   = -1
			relId    = -1
			targetId = -1
		)

		for i, attr := range token.Attr {
			switch attr.Key {
			case "href":
				hrefId = i
			case "rel":
				relId = i
			case "target":
				targetId = i
			}
			if attr.Key == "href" {
				hrefId = i
			}
		}

		if hrefId == -1 {
			return
		}
		if strings.HasPrefix(token.Attr[hrefId].Val, "#") {
			return
		}

		if targetId >= 0 {
			token.Attr[targetId].Val += "_blank"
		} else {
			token.Attr = append(token.Attr, html.Attribute{Key: "target", Val: "blank"})
		}

		if relId >= 0 {
			token.Attr[relId].Val += " noopener noreferrer"
		} else {
			token.Attr = append(token.Attr, html.Attribute{Key: "rel", Val: "noopener noreferrer"})
		}
	}

	tokenizer := html.NewTokenizerFragment(src, "body")
	for {
		token := tokenizer.Next()
		switch token {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("failed to parse html: %w", err)
			}
		case html.StartTagToken:
			token := tokenizer.Token()
			if strings.ToLower(token.Data) == "a" {
				updateLink(&token)

				_, err := dst.Write([]byte(token.String()))
				if err != nil {
					return fmt.Errorf("failed to write out modified token: %w", err)
				}
				continue
			}
			fallthrough
		default:
			_, err := dst.Write(tokenizer.Raw())
			if err != nil {
				return fmt.Errorf("failed to write out token: %w", err)
			}
		}
	}
}
