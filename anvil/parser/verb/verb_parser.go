package verb

import (
	"errors"
	"mutably/anvil/model"
	"mutably/anvil/model/inflection"
	"mutably/anvil/parser"
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

// NewVerbParser creates a *VerbParser that is connected to db and
// uses threadCount workers.
//
// threadCount must be > 0.
//
// pageLimit indicates how many pages should be consumed. If set to -1,
// Parse will always return a true value. When set to N, Parse will
// begin returning false after it has been called N times.
// Valid values are: {-1} U [1, INT_MAX]
func NewVerbParser(db model.Database, threadCount, pageLimit int,
	conjugators *inflection.Conjugators) (*VerbParser, error) {
	if threadCount < 1 {
		return nil, errors.New("Thread count for VerbConsumer must be at least 1")
	}
	if pageLimit < -1 || pageLimit == 0 {
		return nil, errors.New("pageLimit must be in {-1} U [1, INT_MAX]")
	}

	// TODO: Find optimal job queue size.
	const jobQueueSize = 10000

	vparser := &VerbParser{
		PageLimit:     pageLimit,
		PagesConsumed: 0,
		jobQueue:      make(chan parser.Page, jobQueueSize),
	}
	vparser.spawnWorkers(threadCount, db, conjugators)

	return vparser, nil
}

// spawnWorkers creates workerCount parallel workers that take jobs from the
// job queue and store results in db.
func (vparser *VerbParser) spawnWorkers(workerCount int, db model.Database,
	conjugators *inflection.Conjugators) {
	vparser.waitGroup.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func(wkr worker, wg *sync.WaitGroup) {
			defer wg.Done()
			wkr.Start()
		}(NewWorker(db, vparser.jobQueue, conjugators), &vparser.waitGroup)
	}
}

// Wait requests that all workers finish processing page content.
// You should call this method before interacting with results produced
// by this parser.
func (vparser *VerbParser) Wait() {
	close(vparser.jobQueue)
	vparser.waitGroup.Wait()
}

// Parse searches page for verbs and adds their templates to a database.
//
// These templates explain what form the verb is in (e.g., infinitive or
// indicative) and how it can change depending on the context of its
// usage (e.g., 1st person, singular, present tense).
//
// This is a mostly nonblocking call. You should invoke Wait to ensure
// results ready.
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
