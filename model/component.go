package model

import (
	"database/sql"
	"errors"
)

type Language string

// LanguageExists checks db for the existence of lang.
func (db *PsqlDB) LanguageExists(lang Language) (exists bool) {
	db.QueryRow(
		`
		SELECT EXISTS(
			SELECT * FROM languages WHERE description = $1
		)
		`, lang).Scan(&exists)
	return exists
}

type Verb struct {
	Language Language
	Text     string
}

// InsertVerb adds verb to db if it was not already there.
// The id of the new (or existing) verb is returned.
// An error is returned if the verb wasn't present but could not be inserted;
// otherwise, that value is nil.
func (db *PsqlDB) InsertVerb(verb Verb) (int, error) {
	verbId := -1

	row := db.QueryRow(
		`
		SELECT id FROM verbs WHERE verb = $1 AND lang = $2
		`, verb.Text, verb.Language)

	if row.Scan(&verbId) == sql.ErrNoRows {
		row = db.QueryRow(
			`
			INSERT INTO verbs (verb, lang)
				VALUES ($1, $2)
			RETURNING id
			`, verb.Text, verb.Language)
		if row.Scan(&verbId) == sql.ErrNoRows {
			return verbId, errors.New("Failed to insert verb " + verb.Text)
		}
	}
	return verbId, nil
}

type VerbTemplate string

// InsertTemplate adds template to db.
// verbId should be the id of the verb that the template represents. If you
// are unsure of how to get this value, see PsqlDB.InsertVerb.
func (db *PsqlDB) InsertTemplate(template VerbTemplate, verbId int) error {
	_, err := db.Exec(
		`
		INSERT INTO templates (verb_id, template)
			VALUES($1, $2)
		`, verbId, template)

	return err
}
