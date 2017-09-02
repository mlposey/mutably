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

	// The section of the current language being processed
	// This will change throughout the lifetime of Consume because
	// each page has many sections.
	CurrentSection string

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
		templatePattern: regexp.MustCompile(`{{2}[^{]*verb[^{]*}{2}`),
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
		consumer.CurrentSection = (*content)[languageHeaders[i][1]:languageHeaders[i + 1][0]]

		if strings.Contains(consumer.CurrentSection, "===Verb===") {
			consumer.VerbCount++

			language := extractLanguage(content, languageHeaders[i])

			// TODO: Insert new languages into DB.
			// The problem is that we know the description (i.e., the language
			// var itself) but not the tag. Either (a) create some value
			// to store in the tag column or (b) retrieve tags from the web.
			if !language.ExistsIn(consumer.DB) {
				fmt.Println("Language", language, "is undefined")
				continue
			}

			verbTemplates := consumer.GetTemplates(&page.Title, &language)

			for _, template := range verbTemplates {
				err := template.AddTo(consumer.DB)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
	return true, nil
}

// GetTemplates creates a *Verb for each unique verb template in the current section.
func (consumer *VerbConsumer) GetTemplates(verb *string, language *Language) []*Verb {
	templates := consumer.templatePattern.FindAllString(
		consumer.CurrentSection, -1)

	var verbs []*Verb
	haveSeen := make(map[string]bool)

	for _, template := range templates {
		if !haveSeen[template] {
			verbs = append(verbs, &Verb{Text:verb, Lang:language, Template:template})
			haveSeen[template] = true
		}
	}

	return verbs
}

// Find the language header in str.
func extractLanguage(str *string, indices []int) Language {
	return Language(strings.ToLower((*str)[indices[0]+2 : indices[1]-3]))
}
