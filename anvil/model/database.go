package model

// A Database handles queries to a collection of application data.
type Database interface {
	InsertLanguage(*Language)
	InsertWord(string) (wordId int)
	InsertVerb(wordId int, languageId int, tableId int) (verbId int, err error)

	// TODO: Shrink Database interface.
	//       A new interface could work solely with inflection tables. It could
	//       allow mass insertion according to a template or single-column
	//       updates for tables that already exist.
	InsertInfinitive(word string, languageId int) (tableId int)
	InsertAsTense(verb *Verb, tense, person string, isPlural bool) error
	InsertPlural(word, tense string, tableId int) error
}

// KeyRing contains credentials for connecting to a database.
type KeyRing struct {
	DatabaseName string
	Host         string
	Port         uint
	User         string
	Password     string
}
