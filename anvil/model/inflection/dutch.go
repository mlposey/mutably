package inflection

import (
	"errors"
	"log"
	"mutably/anvil/model"
	"regexp"
	"sync"
)

// Dutch is a conjugator for Dutch verbs.
type Dutch struct {
	language *model.Language
	database model.Database

	person *regexp.Regexp
	mood   *regexp.Regexp
	number *regexp.Regexp
	tense  *regexp.Regexp
	infRef *regexp.Regexp

	// A cache of infinitive verbs and their inflection table ids
	tableCache cache
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
		infRef:     regexp.MustCompile(`\|([^\|]+)}}`),
		tableCache: cache{m: make(map[string]int)},
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
	// It is an infinitive.
	if verb.Template == "{{nl-verb}}" {
		verb.TableId = dutch.database.InsertInfinitive(verb.Text,
			dutch.GetLanguage().Id)

		dutch.tableCache.Lock()
		dutch.tableCache.m[verb.Text] = verb.TableId
		dutch.tableCache.Unlock()

		// Note: Past tense plurals are verb forms and won't get caught here.
		err := dutch.database.InsertPlural(verb.Text, "present", verb.TableId)
		return err
	}

	// It is a verb-form.
	infinitive := dutch.infRef.FindStringSubmatch(verb.Template)
	if infinitive == nil {
		return errors.New("Invalid template for verb " + verb.Text)
	}

	dutch.tableCache.RLock()
	if tableId, ok := dutch.tableCache.m[infinitive[1]]; ok {
		dutch.tableCache.RUnlock()
		verb.TableId = tableId
	} else {
		dutch.tableCache.RUnlock()
		verb.TableId = dutch.database.InsertInfinitive(infinitive[1],
			verb.LanguageId)

		dutch.tableCache.Lock()
		dutch.tableCache.m[infinitive[1]] = verb.TableId
		dutch.tableCache.Unlock()
	}

	dutch.addToTable(verb)
	return nil
}

// TODO: Separate this method into other, smaller ones.
func (dutch *Dutch) addToTable(verb *model.Verb) {
	moods := dutch.mood.FindStringSubmatch(verb.Template)
	if moods != nil && moods[1] != "ind" {
		// Some finite verb templates don't display a mood, and all infinitves
		// don't. That means it's important that only finite verbs enter
		// this method.
		return
	}

	tenses := dutch.tense.FindStringSubmatch(verb.Template)
	if tenses == nil {
		return
	}
	var tense string
	switch tenses[1] {
	case "pres":
		tense = "present"
	case "past":
		tense = "past"
	default:
		return
	}

	numbers := dutch.number.FindStringSubmatch(verb.Template)
	if numbers == nil {
		return
	}
	if numbers[1] == "pl" {
		err := dutch.database.InsertAsTense(verb, tense, "", true)
		if err != nil {
			log.Println(err)
		}
	} else {
		persons := dutch.person.FindStringSubmatch(verb.Template)
		if persons == nil {
			// No explicit person means the verb exists for all three.
			persons = []string{"", "123"}
		}
		var person string
		for _, p := range persons[1] {
			switch p {
			case '1':
				person = "first"
			case '2':
				person = "second"
			case '3':
				person = "third"
			}
			err := dutch.database.InsertAsTense(verb, tense, person, false)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// Used by Dutch for caching table ids
type cache struct {
	m map[string]int
	sync.RWMutex
}
