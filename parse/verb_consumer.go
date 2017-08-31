package parse

import (
	"strings"
	"regexp"
	"database/sql"
)

// VerbConsumer adds verbs to a database.
// See *VerbConsumer.Consume and parse.ProcessPages for context.
type VerbConsumer struct {
	DB *sql.DB  // Database connection

	// Pattern for language headers
	languagePattern *regexp.Regexp
	// Pattern for verb templates
	templatePattern *regexp.Regexp

	// These are really only useful for testing.
	// TODO: Test section and verb detection without these variables.
	LanguageCount int
	VerbCount int
}

func NewVerbConsumer(db *sql.DB) *VerbConsumer {
	return &VerbConsumer{
		DB: db,
		languagePattern: regexp.MustCompile(`(==|^==)([\w ]+)(==$|==\s)`),
		templatePattern: regexp.MustCompile(`{{.*}}`),
	}
}

// Consume conditionally adds the contents of page to a database.
//
// If any section of page contains a verb definition, that verb and its
// metadata are inserted in the database defined by consumer.Key. Pages
// that do not contain verb definitions are ignored.
func (consumer *VerbConsumer) Consume(page Page) (bool, error) {
	content := &page.Revision.Text

	// The contents of the page are split into sections that each start
	// with a language header.
	// Each section describes the word as it exists in that language.
	// An English section would start with a '==English==' header.
	languageHeaders := consumer.languagePattern.FindAllStringIndex(*content, -1)

	sectionCount := len(languageHeaders)
	consumer.LanguageCount = sectionCount

	// Section content exists between two language headers. Thus, we must
	// create a fake header at the end to grab the last section.
	languageHeaders = append(languageHeaders, []int{len(*content), 0})

	// TODO: Submit batches of verbs to a multithreaded database worker.
	for i := 0; i < sectionCount - 1; i++ {
		section := (*content)[languageHeaders[i][1]:languageHeaders[i + 1][0]]

		if hasVerb(content, i, languageHeaders) {
			consumer.VerbCount++

			language := extractLanguage(content, languageHeaders[i])
			verbs := consumer.GetTemplates(&section, &page.Title, &language)
			for _, verb := range verbs {
				verb.AddTo(consumer.DB)
			}
		}
	}
	return true, nil
}

// GetTemplates creates a *Verb for each context a language defines.
// For example, the verb 'lie' in English can mean different things, so
// a template for each meaning is assigned to a *Verb object.
//
// sectionBounds should define the start and stop positions within
// pageContent that define verb in language.
func (consumer *VerbConsumer) GetTemplates(section, verb, language *string) []*Verb {
	templates := consumer.templatePattern.FindAllString(*section, -1)

	var verbs []*Verb
	for _, template := range templates {
		verbs = append(verbs, &Verb{Text:verb, Language:language, Template:template})
	}

	return verbs
}

// Find the language header in str.
func extractLanguage(str *string, indices []int) string {
	return (*str)[indices[0]+2 : indices[1]-3]
}

// Return true if the ith section in str contains a verb definition.
func hasVerb(str *string, i int, indices [][]int) bool {
	return strings.Contains((*str)[indices[i][1]:indices[i+1][0]],
		"===Verb===")
}


