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
	insertSingularVerb *sql.Stmt
	insertPluralVerb   *sql.Stmt
}

// NewPsqlDB creates a *PsqlDB using keyring for credentials.
func NewPsqlDB(key KeyRing) (*PsqlDB, error) {
	cred := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		key.User, key.Password, key.Host, key.Port, key.DatabaseName)

	db, err := sql.Open("postgres", cred)
	if err != nil {
		return nil, err
	}

	// TODO: Replace with sql.PingContextCall.
	// Wait for the connection to go through.
	isConnected, remainingTries := false, 10
	for !isConnected && remainingTries > 0 {
		time.Sleep(time.Second * 1)
		remainingTries--
		isConnected = db.Ping() == nil
	}

	if !isConnected {
		db.Close()
		return nil, errors.New("Failed to establish database connection")
	}

	psqlDB := &PsqlDB{db, nil, nil}
	return psqlDB, psqlDB.prepareStatments()
}

func (db *PsqlDB) prepareStatments() error {
	var err error
	db.insertPluralVerb, err = db.Prepare(`
		INSERT INTO verb_forms (lang_id, word_id, inf_id, tense_id, num)
		VALUES
		  ($1,
		  (SELECT id FROM words WHERE word = $2),
		  $3,
		  (SELECT id FROM tenses WHERE tense = $4),
		  2)
	`)
	if err != nil {
		return err
	}
	db.insertSingularVerb, err = db.Prepare(`
		INSERT INTO verb_forms (lang_id, word_id, inf_id, tense_id, person, num)
		VALUES
		  ($1,
		  (SELECT id FROM words WHERE word = $2),
		  $3,
		  (SELECT id FROM tenses WHERE tense = $4),
		  $5,
		  1)
	`)
	if err != nil {
		return err
	}
	return nil
}

// InsertLanguage adds language to the database and sets its Id field.
// If language already exists, the insertion is skipped.
func (db *PsqlDB) InsertLanguage(language *Language) error {
	err := db.QueryRow(`
		SELECT id FROM languages
		WHERE name = $1`,
		language.String(),
	).Scan(&language.Id)

	if err == nil {
		return nil
	} else if err == sql.ErrNoRows {
		err = db.QueryRow(`
			INSERT INTO languages (name)
			VALUES ($1) RETURNING id`,
			language.String(),
		).Scan(&language.Id)
	}
	return err
}

// InsertWord adds word to db if it was not already there.
// The id of the new (or existing) word is returned.
func (db *PsqlDB) InsertWord(word string) int {
	var wordId int
	row := db.QueryRow(`SELECT id FROM words WHERE word = $1`, word)
	if row.Scan(&wordId) == sql.ErrNoRows {
		db.QueryRow(`
			INSERT INTO words (word)
			VALUES ($1) RETURNING id`,
			word,
		).Scan(&wordId)
	}
	return wordId
}

// InsertVerbForm adds a verb form to the database.
// verb should have all fields (except maybe Person) populated.
func (db *PsqlDB) InsertVerbForm(verb *VerbForm) error {
	var tense string
	if verb.Tense == Present {
		tense = "present"
	} else if verb.Tense == Past {
		tense = "past"
	}

	var err error
	if verb.Number == Singular {
		_, err = db.insertSingularVerb.Exec(verb.LanguageId, verb.Word,
			verb.InfinitiveId, tense, verb.Person)
	} else {
		_, err = db.insertPluralVerb.Exec(verb.LanguageId, verb.Word,
			verb.InfinitiveId, tense)
	}
	return err
}
