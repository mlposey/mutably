package model

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"time"
	"errors"
)

// A KeyRing contains credentials for connecting to a database.
type KeyRing struct {
	Database string
	Host string
	Port uint
	User string
	Password string
}

// Validate ensures the database credentials in ring are valid and the target
// is accepting connections. If the key ring is valid, a pointer to an open
// database is returned.
func (ring KeyRing) Validate() (*sql.DB, error) {
	cred := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		ring.User, ring.Password, ring.Host, ring.Port, ring.Database)

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
	return psqlDB, nil
}
