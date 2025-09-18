//spellchecker:words generator
package generator

//spellchecker:words context errors slog runtime sync time
import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"time"
)

// Generator is a dead simple static file generator.
type Generator struct {
	// Inputs are inputs to the generator.
	Inputs []*Scanner

	// Indexes are special templates which are passed all previously generated files.
	Indexes []IndexTemplate

	// ContentTemplate is the template applied to all non-raw files.
	ContentTemplate ContentTemplate

	// PostProcessors are applied to all files.
	// They are applied in order to each file being output.
	PostProcessors []PostProcessor

	// Output is used to write output files.
	Output FileWriter
}

var errRecursiveIndex = errors.New("indexer produced file to be indexed: not allowed")

// debounce debounces the given channel under the given duration.
// A struct received on c is only passed through to the output channel if no other signal arrives within the given window.
func debounce[T any](c <-chan T, window time.Duration, returnLast bool) <-chan T {
	if c == nil {
		return nil
	}
	d := make(chan T)
	go func() {
		defer close(d)

		timer := time.NewTimer(window)
		timer.Stop()

		for v := range c {
		inner:
			for {
				timer.Reset(window)
				select {
				case last := <-c:
					if returnLast {
						v = last
					}
				case <-timer.C:
					d <- v
					break inner
				}
			}

		}
	}()

	return d

}

// Run runs the static site generator with the given context, logging to the given logger.
//
// If context is nil, uses a background context instead.
// If logger is nil, discards all output.
func (generator *Generator) Run(ctx context.Context, logger *slog.Logger) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	ourContext, cancel := context.WithCancel(ctx)
	errChan := make(chan error, 1)

	// registerError registers an error and cancels the context
	registerError := func(err error) {
		select {
		case errChan <- err:
			cancel()
		default:
		}
	}

	var bufferSize = runtime.NumCPU()

	var (
		inputs         = make(chan ScannedFile, bufferSize) // inputs from the inputs
		inputProducers sync.WaitGroup                       // waits for anything writing to inputs

		contents         = make(chan FileWithMetadata, bufferSize) // non-raw files to be wrapped in ContentTemplate
		contentProducers sync.WaitGroup                            // waits for anything writing to the contents

		index          = make(chan ScannedFile, bufferSize)
		indexProducers sync.WaitGroup // anything producing an index entry

		posts         = make(chan File, bufferSize) // outputs to be post-processed
		postProducers sync.WaitGroup                // waits for anything producing post-processing output

		outputs         = make(chan File, bufferSize) // final outputs
		outputProducers sync.WaitGroup                // anything producing final output

		fileWriters sync.WaitGroup
	)

	// start all the inputs
	for i, scanner := range generator.Inputs {
		inputProducers.Add(1)
		go func() {
			defer inputProducers.Done()

			if err := scanner.scan(ourContext, logger, inputs); err != nil {
				registerError(fmt.Errorf("scanner %d failed to scan: %w", i, err))
			}
		}()
	}

	// collect inputs and send them to the appropriate next stages
	postProducers.Add(1)
	contentProducers.Add(1)
	indexProducers.Add(1)
	go func() {
		defer postProducers.Done()
		defer contentProducers.Done()
		defer indexProducers.Done()

		var indexed []IndexEntry
		for result := range inputs {
			if result.Raw {
				posts <- result.File
			} else {
				contents <- result.FileWithMetadata
			}

			if result.Indexed {
				indexed = append(indexed, IndexEntry{Path: result.Path, Metadata: result.Metadata})
			}
		}

		if err := generator.renderIndexes(ourContext, logger, indexed, index); err != nil {
			registerError(fmt.Errorf("failed to render indexes: %w", err))
		}
	}()

	// send indexes to the next appropriate stages
	contentProducers.Add(1)
	postProducers.Add(1)
	go func() {
		defer contentProducers.Done()
		defer postProducers.Done()

		for result := range index {
			if result.Indexed {
				registerError(errRecursiveIndex)
				return
			}
			if result.Raw {
				posts <- result.File
			} else {
				contents <- result.FileWithMetadata
			}
		}
	}()

	// clear the output
	if err := generator.Output.Reset(); err != nil {
		logger.Error("resetting output failed", slog.Any("error", err))
		return err
	}

	// renderContent -> postProcess -> output
	pipe(ourContext, logger, posts, contents, &postProducers, registerError, generator.renderContent)
	pipe(ourContext, logger, outputs, posts, &outputProducers, registerError, generator.postProcess)
	drain(ourContext, logger, outputs, &fileWriters, registerError, generator.Output.Write)

	// close all the components once done
	go func() {
		defer close(inputs)
		inputProducers.Wait()
	}()

	go func() {
		defer close(contents)
		contentProducers.Wait()

	}()

	go func() {
		defer close(posts)
		postProducers.Wait()
	}()

	go func() {
		defer close(index)
		indexProducers.Wait()
	}()

	go func() {
		defer close(outputs)
		outputProducers.Wait()
	}()

	// wait for all the files to have been output
	fileWriters.Wait()

	// and show an error, if any
	select {
	case err := <-errChan:
		logger.Error("build process failed", slog.Any("error", err))
		return err
	default:
		return nil
	}
}

// pipe pipes content from the in channel to the out channel using f.
// when an error occurs aborts and calls registerError instead.
// wg is used to keep track of running operations.
func pipe[S, T any](ctx context.Context, logger *slog.Logger, out chan<- T, in <-chan S, wg *sync.WaitGroup, registerError func(error), f func(context.Context, *slog.Logger, S) (T, error)) {
	drain(ctx, logger, in, wg, registerError, func(ctx context.Context, logger *slog.Logger, input S) error {
		output, err := f(ctx, logger, input)
		if err != nil {
			return fmt.Errorf("pipe processing failed: %w", err)
		}

		select {
		case out <- output:
		case <-ctx.Done():
			return fmt.Errorf("pipe output failed: %w", ctx.Err())
		}
		return nil
	})
}

// drain drains the given channel, and processes f.
// when an error occurs aborts and calls registerError instead.
// wg is used to keep track of running operations.
func drain[T any](ctx context.Context, logger *slog.Logger, in <-chan T, wg *sync.WaitGroup, registerError func(error), f func(context.Context, *slog.Logger, T) error) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			for {
				select {
				case input, ok := <-in:
					if !ok {
						return
					}

					wg.Add(1)
					go func() {
						defer wg.Done()

						err := f(ctx, logger, input)
						if err != nil {
							registerError(fmt.Errorf("drain processing failed: %w", err))
						}
					}()

				case <-ctx.Done():
					registerError(fmt.Errorf("drain input failed: %w", ctx.Err()))
					return
				}
			}
		}
	}()
}
