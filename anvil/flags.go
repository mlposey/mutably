package main

import (
	"flag"
	"fmt"
	"os"
)

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

// GetIntent uses command-line flags to decide what the user wants this
// application to do. That intent is wrapped in the returning function.
//
// flags should be populated prior to calling this function. Refer to
// the documentation in the flag package for details on mapping input
// values to the struct.
func (flags *AppFlags) GetIntent() func() {
	command := os.Args[1]
	// flag.Parse will stop parsing if it notices we tried to use a flag
	// without a dash. Quick! Hide the evidence.
	os.Args = append(os.Args[:1], os.Args[2:]...)
	// Nothing to see here...
	flag.Parse()

	switch command {
	case "import":
		return func() { Import(flags) }

	case "view":
		return View

	case "help":
		return flag.PrintDefaults

	default:
		return func() {
			fmt.Println("Unknown command:", flag.Args()[0])
			ShowHelp()
		}
	}
}
