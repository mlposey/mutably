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

// InsertInfinitive adds a new infinitive verb to the database.
// The id of the new verb's conjugation table is returned.
func (db *PsqlDB) InsertInfinitive(wordId int, languageId int) int {
	var table int
	db.QueryRow(`SELECT add_infinitive($1, $2)`, wordId, languageId).Scan(&table)
	return table
}

// GetTableId returns the conjugation table id for the verb identified by
// languageId and word. An error is returned if the verb does not exist.
func (db *PsqlDB) GetTableId(languageId int, word string) (int, error) {
	var tableId int
	err := db.QueryRow(
		`
		SELECT conjugation_table
		FROM   verbs
		WHERE  lang_id = $1
		AND    word_id = (
			   SELECT id FROM words
			   WHERE  word = $2
		)
		`, languageId, word,
	).Scan(&tableId)

	return tableId, err
}
