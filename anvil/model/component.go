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

// Verb defines a word that is a verb in a specific language.
type Verb struct {
	Id         int
	WordId     int
	LanguageId int
	Text       string
	TableId    int
	Template   string
}

// NewVerb creates a new languageless Verb instance ready to pass to a
// Conjugator.
func NewVerb(wordId, languageId int, text, template string) *Verb {
	return &Verb{
		WordId:     wordId,
		LanguageId: languageId,
		Text:       text,
		Template:   template,
	}
}
