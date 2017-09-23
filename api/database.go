package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// A Database facilitates interaction with a collection of data and ensures
// that both input and output are received according to relative interface
// specifications.
//
// If you find yourself using sql.DB directly then you're probably doing
// something wrong.
type Database interface {
}

// PsqlDB implements the Database interface for PostgreSQL.
type PsqlDB struct {
	db *sql.DB
}

// NewDB creates and returns a PsqlDB instance.
// A non-nil error is returned if there was a problem connecting to the
// database.
func NewDB(host, name, user, password string) (*PsqlDB, error) {
	credentials := fmt.Sprintf("dbname=%s user=%s password=%s host=%s",
		name, user, password, host)

	db, err := sql.Open("postgres", credentials)
	if err != nil {
		return nil, err
	}
	return &PsqlDB{db}, db.Ping()
}
