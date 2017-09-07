package verb

import (
	"anvil/model"
	"anvil/parser"
	"github.com/moovweb/rubex"
	"log"
	"strings"
)

// worker processes Pages that it receives from VerbParser.
type worker struct {
	// The application database
	database model.Database

	// The current language section in the Page.
	// Pages can have multiple sections that the Worker will need to
	// process. This is the one being worked on now.
	languageSection string

	// Pattern for language headers
	languagePattern *rubex.Regexp
	// Pattern for verb templates
	templatePattern *rubex.Regexp

	JobQueue chan parser.Page
}

// NewWorker creates a worker ready to accept jobs from jobPool.
func NewWorker(db model.Database, jobQueue chan parser.Page) worker {
	return worker{
		database:        db,
		languagePattern: rubex.MustCompile(`(?m)^==[^=]+==\n`),
		templatePattern: rubex.MustCompile(`(?m)^{{2}[^{]+verb[^{]+}{2}$`),
		JobQueue:        jobQueue,
	}
}

// Start makes worker begin waiting for jobs from the job pool.
func (wkr worker) Start() {
	for page := range wkr.JobQueue {
		wkr.process(page)
	}
}

func (wkr worker) process(page parser.Page) {
	content := &page.Revision.Text

	// The contents of the page are split into sections that each start
	// with a language header.
	// Each section describes the word as it exists in that language.
	// An English section would start with a '==English==' header.
	languageHeaders := wkr.languagePattern.FindAllStringIndex(*content, -1)

	sectionCount := len(languageHeaders)

	// Section content exists between two language headers. Thus, we must
	// create a fake header at the end to grab the last section.
	languageHeaders = append(languageHeaders, []int{len(*content), 0})

	for i := 0; i < sectionCount; i++ {
		wkr.languageSection =
			(*content)[languageHeaders[i][1]:languageHeaders[i+1][0]]

		if strings.Contains(wkr.languageSection, "===Verb===") {
			verb := model.Verb{
				Language: extractLanguage(content, languageHeaders[i]),
				Text:     page.Title,
			}

			// TODO: Insert new languages into DB.
			// The problem is that we know the description (i.e., the language
			// var itself) but not the tag. Either (a) create some temporary
			// value to store in the tag column or (b) retrieve tags from the web.
			if !wkr.database.LanguageExists(verb.Language) {
				log.Println("Language", verb.Language, "is undefined")
				continue
			}

			verbId, err := wkr.database.InsertVerb(verb)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			verbTemplates := wkr.GetTemplates()

			for _, template := range verbTemplates {
				err := wkr.database.InsertTemplate(template, verbId)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}
}

// GetTemplates creates a VerbTemplate for each unique verb template in the
// current section.
func (wkr worker) GetTemplates() []model.VerbTemplate {
	templates := wkr.templatePattern.FindAllString(wkr.languageSection, -1)

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
