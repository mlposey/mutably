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

// InsertLanguage adds language to the database and sets its Id field.
// If language already exists, the insertion is skipped.
func (db *PsqlDB) InsertLanguage(language *Language) {
	row := db.QueryRow(
		`
		SELECT id
		FROM languages
		WHERE language = $1
		`, language.String(),
	)
	if row.Scan(&language.Id) != sql.ErrNoRows {
		return
	}

	db.QueryRow(
		`
		INSERT INTO languages (language)
		VALUES ($1)
		RETURNING id
		`, language.String(),
	).Scan(&language.Id)
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

// InsertVerb adds a new verb to the database, returning its id.
// An error is returning if the verb could not be inserted.
func (db *PsqlDB) InsertVerb(wordId int, languageId int, tableId int) (int, error) {
	var verbId int
	err := db.QueryRow(
		`
		INSERT INTO verbs (word_id, lang_id, conjugation_table)
		VALUES ($1, $2, $3)
		RETURNING id
		`, wordId, languageId, tableId,
	).Scan(&verbId)
	return verbId, err
}

// InsertInfinitive adds an infinitive verb to the database.
// The id of the new verb's conjugation table is returned. If the verb already
// existed, the id is still returned, but the verb is not inserted.
func (db *PsqlDB) InsertInfinitive(word string, languageId int) int {
	var table int
	// TODO: Return this error.
	db.QueryRow(`SELECT add_infinitive($1, $2)`, word, languageId).Scan(&table)
	return table
}

func (db *PsqlDB) InsertAsTense(verb *Verb, tense, person string,
	isPlural bool) error {
	var tableColumn string
	if isPlural {
		tableColumn = "plural"
	} else {
		tableColumn = person
	}

	// TODO: Find out why parameter substitution doesn't work with Exec.
	//       Note: Query and QueryRow won't commit the results, so those
	//       aren't options.
	_, err := db.Exec(fmt.Sprintf(
		`
		UPDATE tense_inflections
		SET %s = %d
		WHERE id = (
			SELECT %s FROM conjugation_tables
			WHERE id = %d
		)
		`, tableColumn, verb.WordId, tense, verb.TableId,
	))
	return err
}

// InsertPlural adds word to the plural column of a tense inflection table.
func (db *PsqlDB) InsertPlural(word, tense string, conjTableId int) error {
	_, err := db.Exec(fmt.Sprintf(
		`
		UPDATE tense_inflections
		SET plural = (
			SELECT id FROM words
			WHERE word = '%s'
		)
		WHERE id = (
			SELECT %s FROM conjugation_tables
			WHERE id = %d
		)
		`, word, tense, conjTableId,
	))
	return err
}
