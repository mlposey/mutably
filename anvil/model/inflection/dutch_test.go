package inflection_test

import (
	"mutably/anvil/model"
	"mutably/anvil/model/inflection"
	"testing"
)

// Dutch.GetLanguages should return at least one language description.
func TestDutch_hasLanguageDescription(t *testing.T) {
	dutch := inflection.NewDutch()
	if dutch.GetLanguage() == nil {
		t.Error("Conjugators must have at least one language description.")
	}
}

// Dutch.Conjugate should detect if a template defines an infinitive and it
// should add it to the database as both an infinitive and a present plural
// form.
func TestConjugate_detectInfinitive(t *testing.T) {
	db, dutch := makeDutch()
	infinitive := "krijgen"

	err := dutch.Conjugate(infinitive, "{{nl-verb}}")
	check(t, err)

	if db.Words[db.InfinitiveId] != infinitive {
		t.Error("Dutch does not detect pure infinitive templates")
	}
	if db.Plural != infinitive {
		t.Error("Dutch not setting present plural to infinitive")
	}
}

// Dutch.Conjugate should extract the grammatical person from a template
// and convert it to a column name that matches both the context of use
// and naming convention used by the database.
func TestConjugate_identify_person(t *testing.T) {
	db, dutch := makeDutch()
	word := "krijg"

	err := dutch.Conjugate(word, "{{nl-verb form of|p=1|n=sg|t=pres|m=ind|krijgen}}")
	check(t, err)

	if db.First != word || db.Second == word || db.Third == word {
		t.Error("Dutch identified the wrong grammatical person")
	}
}

// Dutch.Conjugate should identify is a verb is of the present or past tense.
func TestConjugate_identify_tense(t *testing.T) {
	db, dutch := makeDutch()

	err := dutch.Conjugate("kreeg", "{{nl-verb form of|n=sg|t=past|m=ind|krijgen}}")
	check(t, err)

	if db.Tense != "past" {
		t.Error("Expected 'past', got", db.Tense)
	}

	err = dutch.Conjugate("krijg", "{{nl-verb form of|p=1|n=sg|t=pres|m=ind|krijgen}}")
	check(t, err)

	if db.Tense != "present" {
		t.Error("Expected 'present', got", db.Tense)
	}
}

// Dutch.Conjugate should only attempt to conjugate the indicative mood for now.
func TestConjugate_check_mood(t *testing.T) {
	db, dutch := makeDutch()
	dutch.Conjugate("krijgt", "{{nl-verb form of|n=pl|m=imp|krijgen}}")

	if db.TableAccessCount != 0 {
		t.Error("Dutch attempted to conjugate mood other than indicative")
	}
}

func makeDutch() (*mockDB, *inflection.Dutch) {
	db := &mockDB{}
	dutch := inflection.NewDutch()
	dutch.SetDatabase(db)
	return db, dutch
}

func check(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}

type mockDB struct {
	Words            []string
	Infinitive       string
	InfinitiveId     int
	First            string
	Tense            string
	Second           string
	Third            string
	Plural           string
	TableAccessCount int
}

func (db *mockDB) InsertLanguage(*model.Language) error { return nil }
func (db *mockDB) InsertWord(word string) (wordId int) {
	db.Words = append(db.Words, word)
	return len(db.Words) - 1
}
func (db *mockDB) InsertVerbForm(verb *model.VerbForm) error {
	db.TableAccessCount++
	db.InfinitiveId = verb.InfinitiveId
	if verb.Tense == model.Present {
		db.Tense = "present"
	} else if verb.Tense == model.Past {
		db.Tense = "past"
	}

	if verb.Number == model.Plural {
		db.Plural = verb.Word
	} else {
		if verb.Person&model.First != 0 {
			db.First = verb.Word
		}
		if verb.Person&model.Second != 0 {
			db.Second = verb.Word
		}
		if verb.Person&model.Third != 0 {
			db.Third = verb.Word
		}
	}
	return nil
}
