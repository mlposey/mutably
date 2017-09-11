package model

// Language is the description of a natural language.
// This includes things like English, Dutch, and German.
type Language string

// Verb defines a word that is a verb in a specific language.
type Verb struct {
	// The language of the verb
	Language Language
	// The actual verb string
	Text string
}

// VerbTemplate describes how a verb presents itself in different contexts.
//
// Templates can be simple (e.g., {{nl-verb}}) or complex (e.g.,
// '{{nl-verb form of|n=sg|t=past|m=subj|treden}}').
// They most often describe some number of grammatical attributes that you
// can read about here: https://en.wikipedia.org/wiki/Finite_verb#Grammatical_categories
type VerbTemplate string
