package inflection

import (
	"errors"
	"log"
	"mutably/anvil/model"
	"regexp"
	"sync"
)

// Dutch is an implementation of Conjugator for the Dutch language.
//
// This struct is guaranteed to be thread safe if instantiated using the
// NewDutch() creational function.
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
		person:   regexp.MustCompile(`p=(\d{1,3})`),
		mood:     regexp.MustCompile(`m=(.{3})`),
		number:   regexp.MustCompile(`n=(.{2})`),
		tense:    regexp.MustCompile(`t=(.{4})`),
		// {{nl-verb form of|...|the_infinitive_reference}}
		infRef:     regexp.MustCompile(`\|([^\|]+)}}`),
		tableCache: cache{m: make(map[string]int)},
	}
}

// GetLanguage returns the English name of the Dutch language.
func (dutch *Dutch) GetLanguage() *model.Language {
	return dutch.language
}

// SetDatabase assigns to dutch a database where it stores results.
func (dutch *Dutch) SetDatabase(db model.Database) error {
	if db == nil {
		return errors.New("Dutch conjugator was given nil database object")
	}
	dutch.database = db
	return nil
}

// Conjugate uses a verb's template to construct parts of the conjugation
// table that it belongs to.
func (dutch *Dutch) Conjugate(verb *model.Verb) error {
	// This header isn't useful for Dutch, but in the future it may be for
	// other languages. Ignore it here rather than in the method that called
	// this one.
	if verb.Template == "{{nl-verb-form}}" {
		return nil
	}

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
		log.Println(*verb)
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

// addToTable adds a verb to the conjugation table of its infinitive form.
// It is important the verb be finite.
func (dutch *Dutch) addToTable(verb *model.Verb) {
	moods := dutch.mood.FindStringSubmatch(verb.Template)
	if moods != nil && moods[1] != "ind" {
		// Some finite verb templates don't display a mood, and all infinitves
		// don't. That means it's important that only finite verbs enter
		// this method.
		return
	}

	tense := dutch.getTense(verb.Template)
	if tense == "" {
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
		for _, person := range dutch.getPersons(verb.Template) {
			err := dutch.database.InsertAsTense(verb, tense, person, false)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// getTense extracts the grammatical tense from a template.
// Returns "" instead of the tense if it is missing.
func (dutch *Dutch) getTense(template string) string {
	tenses := dutch.tense.FindStringSubmatch(template)
	if tenses == nil {
		return ""
	}

	switch tenses[1] {
	case "pres":
		return "present"
	case "past":
		return "past"
	default:
		return ""
	}
}

// getPersons extracts the grammatical person defininitions from a template.
func (dutch *Dutch) getPersons(template string) []string {
	// Grammatical persons (i.e., first, second, third)
	var persons []string

	match := dutch.person.FindStringSubmatch(template)
	if match == nil {
		// No explicit person means the verb exists for all three.
		persons = []string{"first", "second", "third"}
	} else {
		for _, p := range match[1] {
			switch p {
			case '1':
				persons = append(persons, "first")
			case '2':
				persons = append(persons, "second")
			case '3':
				persons = append(persons, "third")
			}
		}
	}
	return persons
}

// Used by Dutch for caching table ids
type cache struct {
	m map[string]int
	sync.RWMutex
}
