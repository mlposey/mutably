package inflection

import (
	"mutably/anvil/model"
)

// A Conjugator builds conjugation tables using verb templates.
// Conjugators should be thread safe.
type Conjugator interface {
	// GetLanguage should return the language that is conjugated.
	GetLanguage() *model.Language
	// SetDatabase should tell the Conjugator where to store results.
	SetDatabase(model.Database) error
	// Conjugate should build part (or all of) a conjugation table.
	Conjugate(verb, template string) error
}
