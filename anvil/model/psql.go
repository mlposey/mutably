package model

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

// PsqlDB implements the Database interface for PostgreSQL.
type PsqlDB struct {
	*sql.DB
}

// NewPsqlDB creates a *PsqlDB using keyring for credentials.
func NewPsqlDB(key KeyRing) (*PsqlDB, error) {
	cred := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		key.User, key.Password, key.Host, key.Port, key.DatabaseName)

	psqlDB, err := sql.Open("postgres", cred)
	if err != nil {
		return nil, err
	}

	// TODO: Replace with sql.PingContextCall.
	// Wait for the connection to go through.
	isConnected, remainingTries := false, 10
	for !isConnected && remainingTries > 0 {
		time.Sleep(time.Second * 1)
		remainingTries--
		isConnected = psqlDB.Ping() == nil
	}

	if !isConnected {
		psqlDB.Close()
		return nil, errors.New("Failed to establish database connection")
	}
	return &PsqlDB{psqlDB}, nil
}

// LanguageExists checks db for the existence of lang.
// If it exists, the id is returned as an integer type.
func (db *PsqlDB) LanguageExists(lang Language) (bool, int) {
	var id int
	row := db.QueryRow(`SELECT id FROM languages WHERE description = $1`, lang)
	if row.Scan(&id) == sql.ErrNoRows {
		return false, -1
	}
	return true, id
}

// InsertWord adds word to db if it was not already there.
// The id of the new (or existing) word is returned.
func (db *PsqlDB) InsertWord(word string) int {
	var wordId int
	row := db.QueryRow(`SELECT id FROM words WHERE word = $1`, word)
	if row.Scan(&wordId) == sql.ErrNoRows {
		db.QueryRow(
			`
			INSERT INTO words (word) VALUES ($1)
			RETURNING id
			`, word,
		).Scan(&wordId)
	}
	return wordId
}

// InsertVerb adds verb to db if it was not already there.
// The id of the new (or existing) verb is returned.
// An error is returned if the verb wasn't present but could not be inserted;
// otherwise, that value is nil.
func (db *PsqlDB) InsertVerb(wordId, languageId int) (int, error) {
	verbId := -1

	row := db.QueryRow(
		`
		INSERT INTO verbs (word_id, lang_id)
		VALUES ($1, $2)
		RETURNING id
		`, wordId, languageId)
	if row.Scan(&verbId) == sql.ErrNoRows {
		return verbId, errors.New("Failed to insert verb")
	}

	return verbId, nil
}

// InsertTemplate adds template to db.
// verbId should be the id of the verb that the template represents. If you
// are unsure of how to get this value, see PsqlDB.InsertVerb.
func (db *PsqlDB) InsertTemplate(template VerbTemplate, verbId int) error {
	var templateId int
	row := db.QueryRow(`SELECT id FROM templates WHERE template = $1`, template)

	if row.Scan(&templateId) == sql.ErrNoRows {
		db.QueryRow(
			`
			INSERT INTO templates (lang_id, template)
			VALUES ((
				SELECT languages.id FROM languages
				JOIN verbs ON verbs.lang_id = languages.id
				WHERE verbs.id = $1
				), $2
			)
			RETURNING id
			`, verbId, template,
		).Scan(&templateId)
	}

	_, err := db.Exec(`
		INSERT INTO verb_templates (verb_id, template_id)
		VALUES ($1,$2)
		`, verbId, templateId)

	return err
}
