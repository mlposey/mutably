package parse

import "database/sql"

type Verb struct {
	Text *string
	Language *string
	Template string
}

// Really, the database should be adding the verb--not the verb adding itself.
// TODO: Solve dependency issues in order to decouple Verb from sql.DB.

func (v *Verb) AddTo(db *sql.DB) error {
	// TODO: Add verb to database.
	return nil
}
