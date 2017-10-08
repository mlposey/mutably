package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

// A Database facilitates interaction with a collection of data and ensures
// that both input and output are received according to relative interface
// specifications.
//
// If you find yourself using sql.DB directly instead of using a Database
// implementation, then you're probably doing something wrong.
type Database interface {
	GetLanguage(int) (*Language, error)
	GetLanguages() ([]*Language, error)
	GetWord(int) (*Word, error)
	GetWords() ([]*Word, error)
	GetUser(string) (*User, error)
	GetUsers() ([]*User, error)
	CreateUser(string, string) (string, error)
	IsAdmin(string) bool
	GetUserId(username, password string) string
	GetConjugationTable(word string) (*ConjugationTable, error)
}

// PsqlDB implements the Database interface for PostgreSQL.
type PsqlDB struct {
	*sql.DB
}

// NewDB creates and returns a PsqlDB instance.
// A non-nil error is returned if there was a problem connecting to the
// database.
func NewDB(host, name, user, password string) (*PsqlDB, error) {
	log.Printf(
		`
-----------------Database Environment-----------------
DATABASE_HOST: %s
DATABASE_NAME: %s
DATABASE_USER: %s
DATABASE_PASSWORD: *redacted*
------------------------------------------------------`,
		host, name, user,
	)
	db, err := sql.Open("postgres", fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s sslmode=disable",
		name, user, password, host,
	))

	if err != nil {
		return nil, err
	}
	return &PsqlDB{db}, db.Ping()
}
