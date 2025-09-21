//spellchecker:words generator
package output

//spellchecker:words context slog http sync testing fstest
import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"testing/fstest"

	"go.tkw01536.de/blog/generator/file"
)

// Server returns a pair of [Output] and [http.Writer].
// The output is intended to be used as the output of a generation, while the handler serves the generated files.
func Server() (Output, http.Handler) {
	dw := &serverWriter{fs: make(fstest.MapFS)}
	return dw, dw.Handler()
}

type serverWriter struct {
	l  sync.Mutex
	fs fstest.MapFS
}

func (sw *serverWriter) Reset() error {
	sw.l.Lock()
	defer sw.l.Unlock()

	clear(sw.fs)
	return nil
}

func (sw *serverWriter) Write(ctx context.Context, logger *slog.Logger, file file.File) error {
	sw.l.Lock()
	defer sw.l.Unlock()

	sw.fs[file.Path] = &fstest.MapFile{
		Data: file.Contents,
	}
	return nil
}

func (sw *serverWriter) Handler() http.Handler {
	return http.FileServerFS(sw.fs)
}
