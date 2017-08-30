package parse

import (
	"strings"
	"regexp"
	"database/sql"
	"fmt"
)

// VerbConsumer adds verbs to a database.
// See *VerbConsumer.Consume and parse.ProcessPages for context.
type VerbConsumer struct {
	DB *sql.DB  // Database connection

	// These are really only useful for testing.
	// TODO: Test section and verb detection without these variables.
	LanguageCount int
	VerbCount int
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

	// Create a placeholder for the last section so we can grab a complete slice.
	// Before:
	// 		"section1 | section2 | section3" --> []{section1, section2}
	// After:
	// 		"section1 | section2 | section3 | placeholder" --> []{section1, section2, section3}
	languageSections = append(languageSections, []int{len(*content), 0})

	for i := 0; i < sectionCount; i++ {
		consumer.tryInsert(*content, page.Title, i, languageSections)
	}
	return true, nil
}

// tryInsert determines if a language context defines a verb. If it does, the
// definition is passed to the insertion procedure.
func (consumer *VerbConsumer) tryInsert(pageContent, word string,
		sectionIndex int, languageSections [][]int) {
	consumer.LanguageCount++
	fmt.Println("language:", extractLanguage(pageContent, languageSections[sectionIndex]))
	if isVerb(pageContent, sectionIndex, languageSections) {
		consumer.VerbCount++
		consumer.insert(
			word,
			extractLanguage(pageContent, languageSections[sectionIndex]),
			"{}", // TODO: Extract the template.
		)
	}
}

// insert adds a verb definition to the database defined by consumer.Key.
func (consumer *VerbConsumer) insert(verb, lang, template string) {
	// TODO: Insert verbs into database.
	fmt.Printf(
`{
	verb: %s,
	lang: %s,
	template: %s
}`, verb, lang, template)
}

// Find the language header in str.
func extractLanguage(str string, indices []int) string {
	return str[indices[0]+2 : indices[1]-3]
}

// Return true if the ith section in str details a verb
func isVerb(str string, i int, indices [][]int) bool {
	return strings.Contains(str[indices[i][1]:indices[i+1][0]],
		"===Verb===")
}


