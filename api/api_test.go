package main_test

import (
	"log"
	"mutably/api"
	"os"
	"testing"
)

var service *main.Service

func init() {
	database, err := main.NewDB(
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
	)
	if err != nil {
		log.Fatal("Could not access database; ", err)
	}

	service, err = main.NewService(database, "8080")
	if err != nil {
		log.Fatal("Could not start service; ", err)
	}
}

func TestMain(m *testing.M) {
	setUp()
	exitCode := m.Run()
	tearDown()
	os.Exit(exitCode)
}

// setUp adds appropriate data to what should be an empty database.
func setUp() {

}

// tearDown reverts all changes made to the database during testing.
func tearDown() {

}
