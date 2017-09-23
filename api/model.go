package main

import (
	"errors"
)

// Language describes a natural language that exists in the database.
type Language struct {
	// The row id of the language
	Id uint `json:"id"`
	// The name of the language (e.g., 'English')
	Name string `json:"name"`
	// The language tag (e.g., 'en')
	Tag string `json:"tag"`
}

// GetLanguage returns from the database the language identified by id.
func (db *PsqlDB) GetLanguage(id int) (*Language, error) {
	return nil, errors.New("GetLanguage is missing implementation")
}

// GetLanguages returns a slice of all languages in the database.
func (db *PsqlDB) GetLanguages() ([]*Language, error) {
	return nil, errors.New("GetLanguages is missing implementation")
}
