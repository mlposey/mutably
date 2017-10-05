package model

import "strings"

// Language is the description of a natural language.
// This includes things like English, Dutch, and German.
type Language struct {
	Id   int
	text string
}

// NewLanguage creates a new Language instance that is not stored in a Database.
func NewLanguage(text string) *Language {
	return &Language{text: text}
}

// String converts language to a canonical form.
func (language *Language) String() string {
	return strings.ToLower(language.text)
}

type GrammaticalTense int

const (
	Present GrammaticalTense = iota
	Past
)

// GrammaticalNumber is a property of finite verbs.
type GrammaticalNumber int

const (
	Singular GrammaticalNumber = 1
	Plural   GrammaticalNumber = 2
)

// GrammaticalPerson is a property of finite verbs.
type GrammaticalPerson int

const (
	First  GrammaticalPerson = 1 << 1
	Second GrammaticalPerson = 1 << 2
	Third  GrammaticalPerson = 1 << 3
)

// VerbForm defines a language-specific form of an infinitive verb.
type VerbForm struct {
	LanguageId   int
	InfinitiveId int

	Word   string
	Tense  GrammaticalTense
	Number GrammaticalNumber
	Person GrammaticalPerson
}
