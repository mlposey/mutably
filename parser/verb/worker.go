package verb

import (
	"anvil/parser"
	"github.com/moovweb/rubex"
)

// Worker handles a unit of VerbParser's work in parallel.
type Worker struct {
	Consumer *VerbParser

	// Pattern for language headers
	languagePattern *rubex.Regexp
	// Pattern for verb templates
	templatePattern *rubex.Regexp

	JobPool chan chan parser.Page
	Job     chan parser.Page
	stop    chan bool
}

// NewWorker creates a worker ready to accept jobs from jobPool.
func NewWorker(consumer *VerbParser, jobPool chan chan parser.Page) Worker {
	return Worker{
		Consumer:        consumer,
		languagePattern: rubex.MustCompile(`(?m)^==[^=]+==\n`),
		templatePattern: rubex.MustCompile(`(?m)^{{2}[^{]+verb[^{]+}{2}$`),
		JobPool:         jobPool,
		Job:             make(chan parser.Page),
		stop:            make(chan bool),
	}
}

// Start makes worker begin waiting for jobs from the job pool.
func (worker Worker) Start() {
	go func() {
		for {
			// Ask the main pool for work.
			worker.JobPool <- worker.Job

			select {
			case page := <-worker.Job:
				worker.Consumer.scrape(page, worker.languagePattern,
					worker.templatePattern)

			case <-worker.stop:
				return
			}
		}
	}()
}

// Stop makes worker stop waiting for jobs from the job pool.
func (worker Worker) Stop() {
	worker.stop <- true
}
