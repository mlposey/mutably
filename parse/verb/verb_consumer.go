package verb

import (
	"anvil/model"
	"anvil/parse"
	"errors"
	"github.com/moovweb/rubex"
	"log"
	"strings"
)

// VerbConsumer adds verbs to a database.
// See *VerbConsumer.Consume and parse.ProcessPages for context.
type VerbConsumer struct {
	// Database connection
	DB model.Database

	// Stop processing pages after Consume is called this many times
	// A value of -1 indicates no limit on the amount of pages consumed.
	PageLimit int
	// The number of pages that have been sent to works. Some may
	// not be fully processed, even if the value has been accounted for.
	PagesConsumed int

	// Workers that process page content in goroutines
	Workers []*Worker
	// A pool of channels for workers that process Pages in separate threads
	WorkerPool chan chan parse.Page
	// This buffered channel is where the jobs will pile up.
	JobQueue chan parse.Page

	// The section of the current language being processed
	// This will change throughout the lifetime of Consume because
	// each page has many sections.
	CurrentSection string
}

// NewVerbConsumer creates a new *VerbConsumer that is connected to db and
// uses threadCount threads.
//
// threadCount must be > 0.
//
// pageLimit indicates how many pages should be consumed. If set to -1,
// Consume will always return a true value. When set to N, Consume will
// begin returning false after it has been called N times.
// Valid values are: {-1} U [1, INT_MAX]
func NewVerbConsumer(db model.Database, threadCount,
	pageLimit int) (*VerbConsumer, error) {
	if threadCount < 1 {
		return nil, errors.New("Thread count for VerbConsumer must be at least 1")
	}
	if pageLimit < -1 {
		return nil, errors.New("pageLimit must be in {-1} U [1, INT_MAX]")
	}

	// Pulled this bad boy out of a hat. Remember: it's not magic
	// if you give it a name ;D
	const queueSize = 5000

	consumer := &VerbConsumer{
		DB:            db,
		PageLimit:     pageLimit,
		PagesConsumed: 0,
		WorkerPool:    make(chan chan parse.Page, threadCount),
		JobQueue:      make(chan parse.Page, queueSize),
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
	for job := range consumer.JobQueue {
		go func(job parse.Page) {
			worker := <-consumer.WorkerPool
			worker <- job
		}(job)
	}
}

// Wait requests that all workers finish processing page content.
func (consumer *VerbConsumer) Wait() {
	for i := range consumer.Workers {
		consumer.Workers[i].Stop()
	}
}

// Parse conditionally adds the contents of page to a database.
//
// If any section of page contains a verb definition, that verb and its
// metadata are inserted in the database defined by consumer.Key. Pages
// that do not contain verb definitions are ignored.
func (consumer *VerbConsumer) Parse(page parse.Page) (bool, error) {
	if consumer.PagesConsumed >= consumer.PageLimit &&
		consumer.PageLimit != -1 {
		return false, errors.New("VerbConsumer is no longer accepting Pages.")
	}
	if consumer.PageLimit != -1 {
		consumer.PagesConsumed++
	}

	// A worker will pick this job up and process page with *VerbConsumer.scrape
	consumer.JobQueue <- page

	return true, nil
}

func (consumer *VerbConsumer) scrape(page parse.Page, languagePattern,
	templatePattern *rubex.Regexp) {
	content := &page.Revision.Text

	// The contents of the page are split into sections that each start
	// with a language header.
	// Each section describes the word as it exists in that language.
	// An English section would start with a '==English==' header.
	languageHeaders := languagePattern.FindAllStringIndex(*content, -1)

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
			if !consumer.DB.LanguageExists(verb.Language) {
				log.Println("Language", verb.Language, "is undefined")
				continue
			}

			verbId, err := consumer.DB.InsertVerb(verb)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			verbTemplates := consumer.GetTemplates(templatePattern)

			for _, template := range verbTemplates {
				err := consumer.DB.InsertTemplate(template, verbId)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}
}

// GetTemplates creates a VerbTemplate for each unique verb template in the
// current section.
func (consumer *VerbConsumer) GetTemplates(p *rubex.Regexp) []model.VerbTemplate {
	templates := p.FindAllString(consumer.CurrentSection, -1)

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
