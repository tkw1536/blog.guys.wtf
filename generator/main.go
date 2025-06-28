package generator

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"sync"
)

// Generator is a dead simple static file generator.
type Generator struct {
	// Inputs are inputs to the generator.
	Inputs []Scanner

	Contents fs.FS // Contents will be processed as markdown.
	Static   fs.FS // Static will be copied as is.

	// Called by index generation to process indexed files.
	IndexCompareFunc IndexComparisonFunc

	// Indexes are templates that are generated after all other files have been processed.
	// The map should map a filename to the template used to render it.
	Indexes map[string]*template.Template

	// ContentTemplate is the Template to use for all final output.
	// It is passed an [HTMLTemplateContext].
	ContentTemplate *Template

	// Output is used to write output files.
	Output FileWriter
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

	var (
		buildWg sync.WaitGroup
		inputWg sync.WaitGroup

		files  = make(chan File)
		inputs = make(chan ScannedFile)
		index  = make(chan Indexed)
	)

	for _, scanner := range generator.Inputs {
		buildWg.Add(1)
		inputWg.Add(1)
		go func() {
			defer buildWg.Done()
			defer inputWg.Done()

			if err := scanner(ctx, logger, inputs); err != nil {
				registerError(fmt.Errorf("scanner %v failed: %w", scanner, err))
			}
		}()
	}

	// close the inputs channel once we're done!
	go func() {
		defer inputWg.Wait()
		close(inputs)
	}()

	var indexedInputs []Indexed

	buildWg.Add(1)
	go func() {
		defer buildWg.Done()

		for input := range inputs {
			if !input.Indexed {
				continue
			}
			indexedInputs = append(indexedInputs, Indexed{
				Path: input.Path,
				Meta: input.Metadata,
			})
		}

	}()

	// build the static directory
	buildWg.Add(1)
	go func() {
		defer buildWg.Done()
		defer close(index)

		if err := generator.scanStatic(ourContext, logger, files); err != nil {
			registerError(fmt.Errorf("failed to scan static directory: %w", err))
		}
	}()

	// build all the markdown
	buildWg.Add(1)
	go func() {
		defer buildWg.Done()

		if err := generator.scanMarkdown(ctx, logger, files, index); err != nil {
			registerError(fmt.Errorf("failed to scan content directory: %w", err))
		}
	}()

	// retrieve and sort the index pages
	indexed := make([]Indexed, 0)

	// scan indexes
	buildWg.Add(1)
	go func() {
		defer buildWg.Done()

		for elem := range index {
			indexed = append(indexed, elem)
		}
	}()

	// render indexes
	go func() {
		buildWg.Wait()
		defer close(files)

		indexed = append(indexed, indexedInputs...)

		if err := generator.RenderIndexes(ctx, logger, indexed, files); err != nil {
			registerError(fmt.Errorf("failed to generate index: %w", err))
			return
		}
	}()

	// write all the output files.
	if err := generator.outputFiles(ctx, logger, files); err != nil {
		registerError(fmt.Errorf("failed to write out files: %w", err))
	}

	select {
	case err := <-errChan:
		logger.Error("build process failed", slog.Any("error", err))
		return err
	default:
		return nil
	}
}
