package verb

import (
	"github.com/moovweb/rubex"
	"log"
	"mutably/anvil/model"
	"mutably/anvil/model/inflection"
	"mutably/anvil/parser"
	"strings"
)

// worker processes Pages that it receives from VerbParser.
type worker struct {
	// The application database
	database model.Database

	// The conjugators that will transform verb templates
	conjugators *inflection.Conjugators

	// The current language section in the Page.
	// Pages can have multiple sections that the Worker will need to
	// process. This is the one being worked on now.
	languageSection string

	// Pattern for matching any header
	headerPattern *rubex.Regexp
	// Pattern for language headers
	languagePattern *rubex.Regexp
	// Pattern for verb headers
	verbPattern *rubex.Regexp
	// Pattern for matching indicative verb templates
	indicativePattern *rubex.Regexp
	// Pattern for matching any verb template
	templatePattern *rubex.Regexp

	jobQueue chan parser.Page
}

// NewWorker creates a worker ready to accept jobs.
func NewWorker(db model.Database, jobQueue chan parser.Page,
	conjugators *inflection.Conjugators) worker {
	return worker{
		database:          db,
		conjugators:       conjugators,
		headerPattern:     rubex.MustCompile(`(?m)^={2,}.*={2,}$`),
		languagePattern:   rubex.MustCompile(`(?m)^==[^=]+==\n`),
		verbPattern:       rubex.MustCompile(`(?m)^={3,}Verb={3,}$`),
		indicativePattern: rubex.MustCompile(`verb( |-)form`),
		templatePattern:   rubex.MustCompile(`(?m)(# )?({{[^{]*}})`),
		jobQueue:          jobQueue,
	}
}

// Start makes worker begin waiting for jobs from the job queue.
func (wkr worker) Start() {
	for page := range wkr.jobQueue {
		wkr.process(page)
	}
}

// process extracts from page the language, word, and verb templates. The
// word and templates are then inserted into wkr.database.
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

	// The id of the word; -1 means it hasn't been inserted into the database.
	wordId := -1

	for i := 0; i < sectionCount; i++ {
		wkr.languageSection =
			(*content)[languageHeaders[i][1]:languageHeaders[i+1][0]]

		if strings.Contains(wkr.languageSection, "===Verb===") {
			verb := model.Verb{
				Language: extractLanguage(content, languageHeaders[i]),
				Text:     page.Title,
			}

			conjugator, err := wkr.conjugators.Get(string(verb.Language))
			if err != nil {
				log.Println(err.Error())
				continue
			}
			// TODO: Find a more efficient place to call this function.
			languageId := wkr.database.InsertLanguage(verb.Language)

			if wordId == -1 {
				wordId = wkr.database.InsertWord(page.Title)
			}

			verbsTemplates := wkr.getTemplates()

			for i := range verbsTemplates {
				verbId, err := wkr.database.InsertVerb(wordId, languageId)
				if err != nil {
					log.Println(err.Error())
					break
				}
				for _, template := range verbsTemplates[i] {
					err := wkr.database.InsertTemplate(template, verbId)
					if err != nil {
						log.Println(err.Error())
					} else {
						conjugator.Conjugate(template)
					}
				}
			}
		}
	}
}

// getTemplates creates a VerbTemplate for each verb template
// in languageSection.
//
// A template can state the form of a verb, offer guidelines on its
// conjugation, or describe the context it is used in. Here are a few
// examples:
//
// {{en-verb}}
// {{en-verb|lies|lying|lied}}
// {{inflection of|lier||3|s|pres|subj|lang=fr}}
func (wkr worker) getTemplates() (templates [][]model.VerbTemplate) {
	verbSections := wkr.getVerbSections()
	if len(verbSections) == 0 {
		return
	}

	templates = make([][]model.VerbTemplate, len(verbSections))

	for vt, verbSection := range verbSections {
		rawTemps := wkr.templatePattern.FindAllStringSubmatch(verbSection, -1)
		if rawTemps == nil {
			break
		}

		base := rawTemps[0][0]
		templates[vt] = append(templates[vt], model.VerbTemplate(base))

		if wkr.indicativePattern.MatchString(base) {
			// This is an indicative verb. It can serve as the template
			// for many tenses and contexts, so there may more template
			// definitions than just the base.
			for i := 1; i < len(rawTemps); i++ {
				var template string
				if len(rawTemps[i]) == 3 {
					// ['# {{atemplate}}', '# ', '{{atemplate}}']
					template = rawTemps[i][2]
				} else {
					// ['{{atemplate}}', '{{atemplate}}']
					template = rawTemps[i][0]
				}
				templates[vt] = append(templates[vt],
					model.VerbTemplate(template))
			}
		}
	}
	return
}

// getVerbSections finds blocks of text within the languageSection that
// detail a verb. These usually begin with a verb header and end at
// either the end of the string or the beginning of a new header.
func (wkr worker) getVerbSections() (sections []string) {
	start, end := 0, 0
	var tmp []int
	for {
		// Find start of section.
		tmp = wkr.verbPattern.FindStringIndex(wkr.languageSection[end:])
		if tmp == nil {
			break
		} else {
			start = end + tmp[1]
		}
		// Find end of section.
		tmp = wkr.headerPattern.FindStringIndex(wkr.languageSection[start:])
		if tmp == nil {
			sections = append(sections, wkr.languageSection[start:])
			break
		} else {
			end = start + tmp[0]
			sections = append(sections, wkr.languageSection[start:end])
		}
	}
	return
}

// Find the language header in str.
func extractLanguage(str *string, indices []int) model.Language {
	return model.Language(strings.ToLower((*str)[indices[0]+2 : indices[1]-3]))
}
