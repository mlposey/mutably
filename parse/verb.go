package parse

import (
	"database/sql"
	"errors"
)

type Verb struct {
	Text *string
	Language *string
	Template string
}

// Really, the database should be adding the verb--not the verb adding itself.
// TODO: Solve dependency issues in order to decouple Verb from sql.DB.

func (v *Verb) AddTo(db *sql.DB) error {
	var verbId uint

	row := db.QueryRow(`SELECT id FROM verbs WHERE verb = $1`, *v.Text)
	if row.Scan(&verbId) == sql.ErrNoRows {
		// Insert the verb so the template has something to refer to.
		row = db.QueryRow(
			`
			INSERT INTO verbs (verb, lang)
				VALUES ($1, $2)
			RETURNING id
			`, *v.Text, *v.Language)
		if row.Scan(&verbId) == sql.ErrNoRows {
			return errors.New("Failed to insert verb " + *v.Text)
		}
	}

	// TODO: Fix duplicate pkey errors.
	// This problem is likely related to the incorrect template extraction.
	// Try to fix that first.
	_, err := db.Exec(
		`
		INSERT INTO templates (verb_id, template)
			VALUES($1, $2)
		`, verbId, v.Template)

	return err
}
