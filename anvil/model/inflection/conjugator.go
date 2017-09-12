package inflection

import (
	"mutably/anvil/model"
	"sync"
)

// A Conjugator builds conjugation tables using verb templates.
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

// conjugators is a singleton that uses the description of a language
// return its conjugator.
// This object can be retrieved by calling GetConjugators().
type conjugators struct {
	c map[string]Conjugator
}

// Get retrieves a Conjugator with a matching language description.
func (conj *conjugators) Get(language string) Conjugator {
	return conj.c[language]
}

var (
	once                sync.Once
	conjugatorsInstance *conjugators
)

// GetConjugators returns the single instance of Conjugators.
func GetConjugators() *conjugators {
	once.Do(func() {
		conjugatorsInstance = &conjugators{make(map[string]Conjugator)}

		items := []Conjugator{
			// Put new Conjugator implementations here.
			&Dutch{},
			// ----------------------------------------
		}
		for _, conjugator := range items {
			for _, language := range conjugator.GetLanguages() {
				conjugatorsInstance.c[string(language)] = conjugator
			}
		}
	})
	return conjugatorsInstance
}
