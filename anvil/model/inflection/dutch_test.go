package inflection_test

import (
	"mutably/anvil/model/inflection"
	"testing"
)

// Dutch.GetLanguages should return at least one language description.
func TestDutch_hasLanguageDescription(t *testing.T) {
	d := &inflection.Dutch{}
	if d.GetLanguage() == nil {
		t.Error("Conjugators must have at least one language description.")
	}
}
