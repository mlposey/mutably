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

	// A cache of word ids for infinitive verbs
	idCache cache
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
		infRef:  regexp.MustCompile(`\|([^\|]+)}}`),
		idCache: cache{m: make(map[string]int)},
	}
}

// GetLanguage returns the English name of the Dutch language.
func (dutch *Dutch) GetLanguage() *model.Language {
	return dutch.language
}

// SetDatabase assigns to dutch a non-nil database where it stores results.
func (dutch *Dutch) SetDatabase(db model.Database) error {
	if db == nil {
		return errors.New("Dutch conjugator was given nil database object")
	}
	dutch.database = db
	return nil
}

// Conjugate uses a verb's template to construct parts of the conjugation
// table that it belongs to. Finite verbs should exist as a word in the
// database before being passed to this method.
func (dutch *Dutch) Conjugate(verb, template string) error {
	// This header isn't useful for Dutch, but in the future it may be for
	// other languages. Ignore it here rather than in the method that called
	// this one.
	if template == "{{nl-verb-form}}" {
		return nil
	}

	if template == "{{nl-verb}}" {
		dutch.handleInfinitive(verb)
		// The first tense plural should just be the infinitive.
		return dutch.handleFinite(verb, "{{nl-verb form of|n=pl|t=pres|m=ind|"+
			verb+"}}")
	} else {
		return dutch.handleFinite(verb, template)
	}
}

// handleInfinitive manages an infinitive verb.
func (dutch *Dutch) handleInfinitive(verb string) {
	infinitiveId := dutch.database.InsertWord(verb)

	dutch.idCache.Lock()
	dutch.idCache.m[verb] = infinitiveId
	dutch.idCache.Unlock()
}

// handleFinite manages a finite verb form.
func (dutch *Dutch) handleFinite(verb, template string) error {
	dutch.database.InsertWord(verb)

	infinitive := dutch.infRef.FindStringSubmatch(template)
	if infinitive == nil {
		return errors.New("Invalid template for verb " + verb)
	}

	// Finite verbs will refer to the verb they come from.
	var infinitiveId int

	dutch.idCache.RLock()
	if id, ok := dutch.idCache.m[infinitive[1]]; ok {
		dutch.idCache.RUnlock()
		infinitiveId = id
	} else {
		dutch.idCache.RUnlock()
		dutch.handleInfinitive(infinitive[1])
		dutch.idCache.RLock()
		infinitiveId = dutch.idCache.m[infinitive[1]]
		dutch.idCache.RUnlock()
	}

	verbForm, e := dutch.assemble(verb, template, infinitiveId)
	if e != nil {
		return e
	}
	//log.Printf("%+v\n", *verbForm)
	return dutch.database.InsertVerbForm(verbForm)
}

// assemble uses a template to assemble the parts of a VerbForm.
func (dutch *Dutch) assemble(verb, template string,
	infinitiveId int) (*model.VerbForm, error) {
	moods := dutch.mood.FindStringSubmatch(template)
	if moods != nil && moods[1] != "ind" {
		// Some finite verb templates don't display a mood, and all infinitves
		// don't. That means it's important that only finite verbs enter
		// this method.
		return nil, errors.New("Expected ind mood")
	}

	tense, err := dutch.getTense(template)
	if err != nil {
		return nil, err
	}

	numberMatch := dutch.number.FindStringSubmatch(template)
	if numberMatch == nil {
		return nil, errors.New("Template is missing grammatical number")
	}

	number, e := dutch.getNumber(template)
	if e != nil {
		return nil, e
	}

	var person int
	if number == model.Singular {
		person = dutch.getPerson(template)
	}

	return &model.VerbForm{
		LanguageId:   dutch.GetLanguage().Id,
		InfinitiveId: infinitiveId,
		Word:         verb,
		Tense:        tense,
		Number:       number,
		Person:       person,
	}, nil
}

// getNumber extracts the grammatical number from a template.
func (dutch *Dutch) getNumber(template string) (model.GrammaticalNumber, error) {
	numberMatch := dutch.number.FindStringSubmatch(template)
	if numberMatch == nil {
		return 0, errors.New("Template is missing grammatical number")
	}

	if numberMatch[1] == "pl" {
		return model.Plural, nil
	} else {
		return model.Singular, nil
	}
}

// getTense extracts the grammatical tense from a template.
func (dutch *Dutch) getTense(template string) (model.GrammaticalTense, error) {
	tenses := dutch.tense.FindStringSubmatch(template)
	if tenses == nil {
		return 0, errors.New("Missing tense for verb form")
	}

	switch tenses[1] {
	case "pres":
		return model.Present, nil
	case "past":
		return model.Past, nil
	default:
		return 0, errors.New("Invalid tense for verb form")
	}
}

// getPersons extracts the grammatical person defininitions from a template.
func (dutch *Dutch) getPerson(template string) int {
	var person int
	match := dutch.person.FindStringSubmatch(template)

	if match != nil {
		for _, p := range match[1] {
			switch p {
			case '1':
				person |= 1 << 1
			case '2':
				person |= 1 << 2
			case '3':
				person |= 1 << 3
			default:
				log.Println("Invalid person in template", template)
			}
		}
	} else {
		// The absence of a person definition means it applies to all persons.
		person = (1 << 1) | (1 << 2) | (1 << 3)
	}
	return person
}

// Cache for infinitive word ids
type cache struct {
	m map[string]int
	sync.RWMutex
}
