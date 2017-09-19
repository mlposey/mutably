package model

// A Database handles queries to a collection of application data.
type Database interface {
	InsertLanguage(*Language)
	InsertWord(string) (wordId int)
	InsertVerb(wordId int, languageId int, tableId int) (verbId int, err error)
	InsertInfinitive(wordId int, languageId int) (tableId int)
	GetTableId(languageId int, word string) (int, error)
}

// KeyRing contains credentials for connecting to a database.
type KeyRing struct {
	DatabaseName string
	Host         string
	Port         uint
	User         string
	Password     string
}
