package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"sync"

	_ "embed"
)

const (
	staticDirectory  = "static"
	contentDirectory = "content"
	outputDirectory  = "public"
)

//go:embed "templates/index.html"
var indexHTML string

var commonData = map[string]any{
	"BlogTitle": "High on Code!",
}

var markdownTemplate = MustTemplate(indexHTML, "index.html", commonData)

func main() {
	var (
		logger = log.Default()
		ctx    = context.Background()
	)

	var errs = make(chan error, 1)
	registerError := func(err error) {
		select {
		case errs <- err:
		default:
		}
	}

	var (
		buildWg sync.WaitGroup
		files   = make(chan File)
		index   = make(chan Indexed)
	)

	// build the static directory
	buildWg.Add(1)
	go func() {
		defer buildWg.Done()

		if err := scanStatic(ctx, logger, files, staticDirectory); err != nil {
			registerError(fmt.Errorf("failed to scan static directory: %w", err))
		}
	}()

	// build markdown
	buildWg.Add(1)
	go func() {
		defer buildWg.Done()
		defer close(index)

		if err := scanMarkdown(ctx, logger, markdownTemplate, files, index, contentDirectory); err != nil {
			registerError(fmt.Errorf("failed to scan content directory: %w", err))
		}
	}()

	// retrieve and sort indexed pages
	indexed := make([]Indexed, 0)
	buildWg.Add(1)
	go func() {
		defer buildWg.Done()

		for elem := range index {
			indexed = append(indexed, elem)
		}

		slices.SortFunc(indexed, func(left, right Indexed) int {
			lDate, _ := left.Meta["date"].(string)
			rDate, _ := right.Meta["date"].(string)

			// sort first by date, then by path
			return cmp.Or(
				strings.Compare(rDate, lDate),
				strings.Compare(left.Path, right.Path),
			)
		})
	}()

	// render the index page, then close the files to signal completion
	go func() {
		buildWg.Wait()
		defer close(files)

		files <- RenderIndex(ctx, logger, indexed, markdownTemplate)
	}()

	// write out all the files
	if err := outputFiles(ctx, logger, outputDirectory, files); err != nil {
		registerError(fmt.Errorf("failed to scan static directory: %w", err))
	}

	select {
	case err := <-errs:
		log.Fatal(err)
	default:

	}
}
