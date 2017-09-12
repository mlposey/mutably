package inflection

import "mutably/anvil/model"

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
