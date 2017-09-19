package inflection

import (
	"log"
	"mutably/anvil/model"
	"regexp"
)

// Dutch is a conjugator for Dutch verbs.
type Dutch struct {
	language *model.Language
	database model.Database

	person  *regexp.Regexp
	mood    *regexp.Regexp
	number  *regexp.Regexp
	tense   *regexp.Regexp
	inf_ref *regexp.Regexp
}

// NewDutch creates and returns a new Dutch instance.
func NewDutch() *Dutch {
	return &Dutch{
		language: model.NewLanguage("Dutch"),
		person:   regexp.MustCompile(`p=(\d{1,2})`),
		mood:     regexp.MustCompile(`m=(.{3})`),
		number:   regexp.MustCompile(`n=(.{2})`),
		tense:    regexp.MustCompile(`t=(.{4})`),
		// {{nl-verb form of|...|the_infinitive_reference}}
		inf_ref: regexp.MustCompile(`\|([^\|]+)}}`),
	}
}

// GetLanguage provides various descriptions of the dutch language.
// For example, "Dutch" is one form. "Nederlands" may be another.
func (dutch *Dutch) GetLanguage() *model.Language {
	return dutch.language
}

// SetDatabase assigns to dutch a database where it stores results.
func (dutch *Dutch) SetDatabase(db model.Database) error {
	dutch.database = db
	return nil
}

// Conjugate uses the template to build part of the conjugation table that
// the word belongs to.
func (dutch *Dutch) Conjugate(verb *model.Verb) error {
	verb.LanguageId = dutch.GetLanguage().Id
	log.Println("Conjugating", *verb)
	return nil
}
