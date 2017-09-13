package main

// AppFlags holds CLI flags passed to the application.
type AppFlags struct {
	// Flags global to all actions

	// The number of pages to process before stopping.
	PageLimit int

	// Toggles verbose logging output
	BeVerbose bool

	// ----------------------------

	// Flags specific to importing

	// The name of a database
	DBName string

	// The name of a user with access to DBName
	DBUser string

	// The password of DBUser
	DBPassword string

	// The host name of DBName's server
	DBHost string

	// The port of DBName's server
	DBPort uint
}
