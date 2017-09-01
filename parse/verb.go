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
	var languageExists bool
	db.QueryRow(
		`
		SELECT EXISTS(
			SELECT * FROM languages WHERE description = $1
		)
		`, *v.Language).Scan(&languageExists)

	// TODO: Handle addition of new languages.
	// Some languages in the wiki won't match the descriptions from a
	// registry file exactly. Create a way to infer that two are the same.
	// Looking at language tags may be a start.
	if !languageExists {
		return errors.New("Language " + *v.Language + " is undefined")
	}

	// TODO: Fix duplicate pkey errors.
	// This problem is likely related to the incorrect template extraction.
	// Try to fix that first.
	_, err := db.Exec(
		`
		INSERT INTO verbs (word, lang, template)
			VALUES($1, $2, $3)
		`, *v.Text, *v.Language, v.Template)

	return err
}
