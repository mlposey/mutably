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
