package inflection

import (
	"mutably/anvil/model"
)

// Dutch is a conjugator for the Dutch language.
type Dutch struct {
}

// GetLanguage provides various descriptions of the dutch language.
// For example, "Dutch" is one form. "Nederlands" may be another.
func (dutch *Dutch) GetLanguages() []model.Language {
	return []model.Language{"Dutch", "Nederlands"}
}

// SetDatabase assigns to dutch a database where it stores results.
func (dutch *Dutch) SetDatabase(model.Database) error {
	return nil
}

// Conjugate uses the template to build part of the conjugation table that
// the word belongs to.
func (dutch *Dutch) Conjugate(model.VerbTemplate) error {
	return nil
}
