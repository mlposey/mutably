package model

import (
	"database/sql"
	"errors"
)

type Language string

// ExistsIn checks db for the presence of the language.
func (lang Language) ExistsIn(db *sql.DB) (exists bool) {
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
	Text string
}

// TryInsert adds verb to db if it is not already there.
// The id of the new (or existing) row is returned.
func (verb Verb) TryInsert(db *sql.DB) (int, error) {
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

// AddTo adds the template of a verb to db.
func (template VerbTemplate) AddTo(db *sql.DB, verbId int) error {
	_, err := db.Exec(
		`
		INSERT INTO templates (verb_id, template)
			VALUES($1, $2)
		`, verbId, template)

	return err
}
