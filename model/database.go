package model

// A Database handles queries to a collection of application data.
type Database interface {
	LanguageExists(Language) bool
	InsertVerb(Verb) (int, error)
	InsertTemplate(template VerbTemplate, verbId int) error
}

// A KeyRing contains credentials for connecting to a database.
type KeyRing struct {
	DatabaseName string
	Host         string
	Port         uint
	User         string
	Password     string
}
