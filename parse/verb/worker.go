package verb

import (
	"anvil/parse"
	"github.com/moovweb/rubex"
)

// A Worker takes care of the verb consumption process for a page.
type Worker struct {
	Consumer *VerbConsumer

	// Pattern for language headers
	languagePattern *rubex.Regexp
	// Pattern for verb templates
	templatePattern *rubex.Regexp

	JobPool chan chan parse.Page
	Job     chan parse.Page
	stop    chan bool
}

// NewWorker creates a worker ready to accept jobs from jobPool.
func NewWorker(consumer *VerbConsumer, jobPool chan chan parse.Page) Worker {
	return Worker{
		Consumer:        consumer,
		languagePattern: rubex.MustCompile(`(?m)^==[^=]+==\n`),
		templatePattern: rubex.MustCompile(`(?m)^{{2}[^{]+verb[^{]+}{2}$`),
		JobPool:         jobPool,
		Job:             make(chan parse.Page),
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
