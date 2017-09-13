package model

// A Database handles queries to a collection of application data.
type Database interface {
	InsertLanguage(Language) (languageId int)
	InsertWord(string) (wordId int)
	InsertVerb(wordId, languageId int) (verbId int, err error)
	InsertTemplate(template VerbTemplate, verbId int) error
}

// KeyRing contains credentials for connecting to a database.
type KeyRing struct {
	DatabaseName string
	Host         string
	Port         uint
	User         string
	Password     string
}
