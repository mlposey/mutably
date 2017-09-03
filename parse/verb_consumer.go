package parse

import (
	"anvil/model"
	"database/sql"
	"errors"
	"log"
	"regexp"
	"strings"
)

// VerbConsumer adds verbs to a database.
// See *VerbConsumer.Consume and parse.ProcessPages for context.
type VerbConsumer struct {
	DB *sql.DB // Database connection

	// Workers that process page content in goroutines
	Workers []*Worker
	// A pool of channels for workers that process Pages in separate threads
	WorkerPool chan chan Page
	// This buffered channel is where the jobs will pile up.
	JobQueue chan Page

	// The section of the current language being processed
	// This will change throughout the lifetime of Consume because
	// each page has many sections.
	CurrentSection string

	// Pattern for language headers
	languagePattern *regexp.Regexp
	// Pattern for verb templates
	templatePattern *regexp.Regexp
}

// NewVerbConsumer creates a new *VerbConsumer that is connected to db and
// uses threadCount threads.
// The connection to db must be valid and threadCount must be > 0
func NewVerbConsumer(db *sql.DB, threadCount int) (*VerbConsumer, error) {
	if threadCount < 1 {
		return nil, errors.New("Thread count for VerbConsumer must be at least 1")
	}

	// Pulled this bad boy out of a hat. Remember: it's not magic
	// if you give it a name ;D
	const queueSize = 1000

	consumer := &VerbConsumer{
		DB:              db,
		WorkerPool:      make(chan chan Page, threadCount),
		JobQueue:        make(chan Page, queueSize),
		languagePattern: regexp.MustCompile(`(==|^==)([\w ]+)(==$|==\s)`),
		templatePattern: regexp.MustCompile(`{{2}[^{]*verb[^{]*}{2}`),
	}

	for i := 0; i < threadCount; i++ {
		worker := NewWorker(consumer, consumer.WorkerPool)
		consumer.Workers = append(consumer.Workers, &worker)
		worker.Start()
	}
	go consumer.coordinateJobs()

	return consumer, nil
}

func (consumer *VerbConsumer) coordinateJobs() {
	for {
		select {
		case job := <-consumer.JobQueue:
			go func(job Page) {
				worker := <-consumer.WorkerPool
				worker <- job
			}(job)
		}
	}
}

// Wait requests that all workers finish processing page content.
func (consumer *VerbConsumer) Wait() {
	for i := range consumer.Workers {
		consumer.Workers[i].Stop()
	}
}

// Consume conditionally adds the contents of page to a database.
//
// If any section of page contains a verb definition, that verb and its
// metadata are inserted in the database defined by consumer.Key. Pages
// that do not contain verb definitions are ignored.
func (consumer *VerbConsumer) Consume(page Page) (bool, error) {

	// A worker will pick this job up and process page with *VerbConsumer.scrape
	consumer.JobQueue <- page

	return true, nil
}

func (consumer *VerbConsumer) scrape(page Page) {
	content := &page.Revision.Text

	// The contents of the page are split into sections that each start
	// with a language header.
	// Each section describes the word as it exists in that language.
	// An English section would start with a '==English==' header.
	languageHeaders := consumer.languagePattern.FindAllStringIndex(*content, -1)

	sectionCount := len(languageHeaders)

	// Section content exists between two language headers. Thus, we must
	// create a fake header at the end to grab the last section.
	languageHeaders = append(languageHeaders, []int{len(*content), 0})

	for i := 0; i < sectionCount-1; i++ {
		consumer.CurrentSection = (*content)[languageHeaders[i][1]:languageHeaders[i+1][0]]

		if strings.Contains(consumer.CurrentSection, "===Verb===") {
			verb := model.Verb{
				Language: extractLanguage(content, languageHeaders[i]),
				Text:     page.Title,
			}

			// TODO: Insert new languages into DB.
			// The problem is that we know the description (i.e., the language
			// var itself) but not the tag. Either (a) create some temporary
			// value to store in the tag column or (b) retrieve tags from the web.
			if !verb.Language.ExistsIn(consumer.DB) {
				log.Println("Language", verb.Language, "is undefined")
				continue
			}

			verbId, err := verb.TryInsert(consumer.DB)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			verbTemplates := consumer.GetTemplates()

			for _, template := range verbTemplates {
				err := template.AddTo(consumer.DB, verbId)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}
}

// GetTemplates creates a VerbTemplate for each unique verb template in the
// current section.
func (consumer *VerbConsumer) GetTemplates() []model.VerbTemplate {
	templates := consumer.templatePattern.FindAllString(
		consumer.CurrentSection, -1)

	var verbTemplates []model.VerbTemplate
	haveSeen := make(map[string]bool)

	// Some sections have repeat template definitions; ignore duplicates.
	for _, template := range templates {
		if !haveSeen[template] {
			verbTemplates = append(verbTemplates, model.VerbTemplate(template))
			haveSeen[template] = true
		}
	}

	return verbTemplates
}

// Find the language header in str.
func extractLanguage(str *string, indices []int) model.Language {
	return model.Language(strings.ToLower((*str)[indices[0]+2 : indices[1]-3]))
}

// ---------------- Multithreading Logic ----------------

// A Worker takes care of the verb consumption process for a page.
type Worker struct {
	Consumer *VerbConsumer
	JobPool  chan chan Page
	Job      chan Page
	stop     chan bool
}

// NewWorker creates a worker ready to accept jobs from jobPool.
func NewWorker(consumer *VerbConsumer, jobPool chan chan Page) Worker {
	return Worker{
		Consumer: consumer,
		JobPool:  jobPool,
		Job:      make(chan Page),
		stop:     make(chan bool),
	}
}

func (worker Worker) Start() {
	go func() {
		for {
			// Ask the main pool for work.
			worker.JobPool <- worker.Job

			select {
			case page := <-worker.Job:
				worker.Consumer.scrape(page)

			case <-worker.stop:
				return
			}
		}
	}()
}

func (worker Worker) Stop() {
	go func() {
		worker.stop <- true
	}()
}
