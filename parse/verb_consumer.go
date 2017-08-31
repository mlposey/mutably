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

	// The contents of the page are split into language sections. Each
	// section describes the word as it exists in that language.
	// An English section would start with a '==English==' header.
	//languageSections := regexp.MustCompile(`(==|^==)(\w+)(==$|==\s)`).
	languageSections := regexp.MustCompile(`(==|^==)([\w ]+)(==$|==\s)`).
		FindAllStringIndex(*content, -1)

	sectionCount := len(languageSections)
	consumer.LanguageCount = sectionCount

	// Create a placeholder for the last section so we can grab a complete slice.
	// Before:
	// 		"section1 | section2 | section3" --> []{section1, section2}
	// After:
	// 		"section1 | section2 | section3 | placeholder" --> []{section1, section2, section3}
	languageSections = append(languageSections, []int{len(*content), 0})

	// TODO: Submit batches of verbs to a multithreaded database worker.
	for i := 0; i < sectionCount - 1; i++ {
		section := (*content)[languageSections[i][1]:languageSections[i + 1][0]]

		if hasVerb(content, i, languageSections) {
			consumer.VerbCount++

			language := extractLanguage(content, languageSections[i])
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


