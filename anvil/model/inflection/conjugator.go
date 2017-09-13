package inflection

import (
	"errors"
	"mutably/anvil/model"
)

// A Conjugator builds conjugation tables using verb templates.
// Conjugators should be thread safe.
type Conjugator interface {
	// GetLanguages should tell us what language we're conjugating.
	// Some languages have many descriptions (e.g., Dutch, Nederlands, Flemish)
	// so they should all be contained in the slice that is returned.
	GetLanguages() []model.Language
	// SetDatabase should tell the Conjugator where to store results.
	SetDatabase(model.Database) error
	// Conjugate should build part (or all of) a conjugation table.
	Conjugate(model.VerbTemplate) error
}

// Conjugators uses the descriptions of languages to return Conjugator structs.
type Conjugators struct {
	c map[string]Conjugator
}

// Makes and returns a new Conjugators instance
func NewConjugators() *Conjugators {
	return &Conjugators{make(map[string]Conjugator)}
}

// Add defines conjugator under each of its language descriptions.
func (conj *Conjugators) Add(conjugator Conjugator) {
	for _, language := range conjugator.GetLanguages() {
		conj.c[string(language)] = conjugator
	}
}

// Get retrieves a Conjugator with a matching language description.
// It will return nil and an error if the language is not defined.
func (conj *Conjugators) Get(language string) (Conjugator, error) {
	if conjugator, exists := conj.c[language]; exists {
		return conjugator, nil
	} else {
		return nil, errors.New("No conjugator exists for language " + language)
	}
}
