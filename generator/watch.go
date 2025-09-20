//spellchecker:words generator
package generator

//spellchecker:words context slog time github farmergreg rfsnotify pkglib errorsx
import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/farmergreg/rfsnotify"
	"go.tkw01536.de/pkglib/errorsx"
	"gopkg.in/fsnotify.v1"
)

// Watch is like calling [Run] every time a signal is received on the given channel.
// No two runs of generator occur simultaneously.
func (generator *Generator) Watch(ctx context.Context, logger *slog.Logger) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	watch, closer, err := generator.newNotifier(ctx, logger)
	if err != nil {
		return fmt.Errorf("failed to watch: %w", err)
	}
	defer closer()

	doBuild := func() {
		logger.Info("triggering rebuild")
		err := generator.Run(ctx, logger)
		if err != nil {
			logger.Error("rebuild failed", slog.Any("err", err))
		}
		logger.Info("rebuild succeeded")
	}

	debouncedSignal := debounce(watch, time.Second, false)

	doBuild()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context closed: %w", ctx.Err())
		case <-debouncedSignal:
			doBuild()
		}

	}
}

// watchScanners watches the inputs of the given scanners.
//
// It returns a triple.
// The first contains a channel that is sent a signal whenever any directory changes.
// The seconds is a function to stop watching and close the channel.
// The third returns an error if watching failed.
func (generator *Generator) newNotifier(ctx context.Context, logger *slog.Logger) (<-chan fsnotify.Event, func() error, error) {
	watcher, err := rfsnotify.NewWatcher()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to watch directories: %w", err)
	}

	c := make(chan fsnotify.Event, 1)
	go func() {
		defer close(c)

		for {
			select {
			case <-ctx.Done():
				return
			case e, ok := <-watcher.Events:
				if !ok {
					return
				}
				logger.Info("watcher triggered", slog.String("event", e.String()))
				c <- e
			case e, ok := <-watcher.Errors:
				logger.Error("watcher saw error", slog.Any("error", e))
				if !ok {
					return
				}
			}
		}
	}()

	for _, scanner := range generator.Inputs {
		for _, dir := range scanner.Paths() {
			if err := watcher.AddRecursive(dir); err != nil {
				err := fmt.Errorf("failed to watch %q: %w", dir, err)
				return nil, nil, errorsx.Combine(err, watcher.Close())
			}
			logger.Info("watching directory", slog.String("directory", dir))
		}

	}

	return c, watcher.Close, nil
}
