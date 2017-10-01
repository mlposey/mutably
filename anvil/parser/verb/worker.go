package verb

import (
	"log"
	"mutably/anvil/model"
	"mutably/anvil/model/inflection"
	"mutably/anvil/parser"
	"regexp"
	"strings"
)

// worker processes Pages that it receives from VerbParser.
type worker struct {
	// The application database
	database model.Database

	// Maps a canonicalized language to its conjugator
	conjugators map[string]inflection.Conjugator

	// The current language section in the Page.
	// Pages can have multiple sections that the Worker will need to
	// process. This is the one being worked on now.
	languageSection string

	// Pattern for matching any header
	headerPattern *regexp.Regexp
	// Pattern for language headers
	languagePattern *regexp.Regexp
	// Pattern for verb headers
	verbPattern *regexp.Regexp
	// Pattern for matching indicative verb templates
	indicativePattern *regexp.Regexp
	// Pattern for matching any verb template
	templatePattern *regexp.Regexp

	jobQueue chan parser.Page
}

// NewWorker creates a worker ready to accept jobs.
func NewWorker(db model.Database, jobQueue chan parser.Page,
	conjugators map[string]inflection.Conjugator) worker {
	return worker{
		database:          db,
		conjugators:       conjugators,
		headerPattern:     regexp.MustCompile(`(?m)^={2,}.*={2,}$`),
		languagePattern:   regexp.MustCompile(`(?m)^==[^=]+==\n`),
		verbPattern:       regexp.MustCompile(`(?m)^={3,}Verb={3,}$`),
		indicativePattern: regexp.MustCompile(`verb( |-)form`),
		templatePattern:   regexp.MustCompile(`(?m)(# )?({{[^{]*}})`),
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

	// Each page defines a word in multiple languages. The definitions are
	// split into sections that begin with a header.
	languageHeaders := wkr.languagePattern.FindAllStringIndex(*content, -1)

	sectionCount := len(languageHeaders)

	// Section content exists between two language headers. So we must
	// create a fake header at the end to grab the last section.
	languageHeaders = append(languageHeaders, []int{len(*content), 0})

	// Pass each verb from the sections to the conjugator with a matching
	// language.
	for i := 0; i < sectionCount; i++ {
		wkr.languageSection =
			(*content)[languageHeaders[i][1]:languageHeaders[i+1][0]]

		if strings.Contains(wkr.languageSection, "===Verb===") {
			language := model.NewLanguage(
				(*content)[languageHeaders[i][0]+2 : languageHeaders[i][1]-3],
			)

			conjugator, ok := wkr.conjugators[language.String()]
			if !ok { // The language isn't supported.
				continue
			}

			templates := wkr.getTemplates()
			for _, template := range templates {
				err := conjugator.Conjugate(page.Title, template)
				if err != nil {
					log.Println(err)
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
func (wkr worker) getTemplates() (templates []string) {
	verbSection := wkr.getVerbSection()
	if verbSection == "" {
		return templates
	}

	rawTemps := wkr.templatePattern.FindAllStringSubmatch(verbSection, -1)
	if rawTemps == nil {
		return templates
	}

	base := rawTemps[0][0]
	templates = append(templates, base)

	if wkr.indicativePattern.MatchString(base) {
		// Indicative verbs can define multiple templates. Get them all.
		for i := 1; i < len(rawTemps); i++ {
			var template string
			if len(rawTemps[i]) == 3 {
				// ['# {{atemplate}}', '# ', '{{atemplate}}']
				template = rawTemps[i][2]
			} else {
				// ['{{atemplate}}', '{{atemplate}}']
				template = rawTemps[i][0]
			}
			templates = append(templates, template)
		}
	}

	return templates
}

// getVerbSections finds blocks of text within the languageSection that
// detail a verb. These usually begin with a verb header and end at
// either the end of the string or the beginning of a new header.
func (wkr worker) getVerbSection() string {
	start, end := 0, 0
	var tmp []int

	// Find start of section.
	tmp = wkr.verbPattern.FindStringIndex(wkr.languageSection[end:])
	if tmp == nil {
		return ""
	} else {
		start = end + tmp[1]
	}
	// Find end of section.
	tmp = wkr.headerPattern.FindStringIndex(wkr.languageSection[start:])
	if tmp == nil {
		return wkr.languageSection[start:]
	} else {
		end = start + tmp[0]
		return wkr.languageSection[start:end]
	}
}
