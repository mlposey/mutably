package verb

import (
	"anvil/model"
	"anvil/parser"
	"errors"
	"sync"
)

// VerbParser uses parallel workers to add verbs to a database.
type VerbParser struct {
	// Stop processing pages after Parse is called this many times
	// A value of -1 indicates no limit on the amount of pages consumed.
	PageLimit int
	// The number of pages that have been sent to works. Some may
	// not be fully processed, even if the value has been accounted for.
	PagesConsumed int

	// This buffered channel holds Page sent from parse.
	jobQueue chan parser.Page
	// A wait group for the workers.
	waitGroup sync.WaitGroup
}

// NewVerbParser creates a new *VerbParser that is connected to db and
// uses threadCount threads.
//
// threadCount must be > 0.
//
// pageLimit indicates how many pages should be consumed. If set to -1,
// Parse will always return a true value. When set to N, Consume will
// begin returning false after it has been called N times.
// Valid values are: {-1} U [1, INT_MAX]
func NewVerbParser(db model.Database, threadCount,
	pageLimit int) (*VerbParser, error) {
	if threadCount < 1 {
		return nil, errors.New("Thread count for VerbConsumer must be at least 1")
	}
	if pageLimit < -1 {
		return nil, errors.New("pageLimit must be in {-1} U [1, INT_MAX]")
	}

	// Pulled this bad boy out of a hat. Remember: it's not magic
	// if you give it a name ;D
	const queueSize = 5000

	vparser := &VerbParser{
		PageLimit:     pageLimit,
		PagesConsumed: 0,
		jobQueue:      make(chan parser.Page, queueSize),
	}

	vparser.waitGroup.Add(threadCount)

	for i := 0; i < threadCount; i++ {
		go func(wkr worker, wg *sync.WaitGroup) {
			defer wg.Done()
			wkr.Start()
		}(NewWorker(db, vparser.jobQueue), &vparser.waitGroup)
	}
	return vparser, nil
}

// Wait requests that all workers finish processing page content.
func (vparser *VerbParser) Wait() {
	close(vparser.jobQueue)
	vparser.waitGroup.Wait()
}

// Parse conditionally adds the contents of page to a database.
//
// If any section of page contains a verb definition, that verb and its
// metadata are inserted in the database defined by consumer.Key. Pages
// that do not contain verb definitions are ignored.
func (vparser *VerbParser) Parse(page parser.Page) (bool, error) {
	if vparser.PagesConsumed >= vparser.PageLimit &&
		vparser.PageLimit != -1 {
		return false, errors.New("VerbParser is no longer accepting Pages.")
	}
	if vparser.PageLimit != -1 {
		vparser.PagesConsumed++
	}

	// Send the page to a worker that waits on the other end.
	vparser.jobQueue <- page

	return true, nil
}
