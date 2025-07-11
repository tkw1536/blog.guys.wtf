package generator

import (
	"context"
	"log/slog"
	"net/http"
	"testing/fstest"
)

// NewDebugServer returns a pair of FileWriter and http.Handler.
// FileWriter can be used as the output of a generator, while the handler serves the generated files.
func NewDebugServer() (FileWriter, http.Handler) {
	dw := &debugWriter{fs: make(fstest.MapFS)}
	return dw.Write, dw.Handler()
}

type debugWriter struct {
	fs fstest.MapFS
}

func (sw *debugWriter) Write(ctx context.Context, logger *slog.Logger, file File) error {
	sw.fs[file.Path] = &fstest.MapFile{
		Data: file.Contents,
	}
	return nil
}

func (sw *debugWriter) Handler() http.Handler {
	return http.FileServerFS(sw.fs)
}
