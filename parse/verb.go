package parse

import "database/sql"

type Verb struct {
	Text *string
	Language *string
	Template string
}

func (v *Verb) AddTo(db *sql.DB) error {
	// TODO: Add verb to database.
	return nil
}
