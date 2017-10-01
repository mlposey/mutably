package model

// A Database handles queries to a collection of application data.
type Database interface {
	InsertLanguage(*Language) error
	InsertWord(string) (wordId int)
	InsertVerbForm(*VerbForm) error
}

// KeyRing contains credentials for connecting to a database.
type KeyRing struct {
	DatabaseName string
	Host         string
	Port         uint
	User         string
	Password     string
}
