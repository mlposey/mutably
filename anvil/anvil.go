package main

import (
	"flag"
	"fmt"
	"log"
	"mutably/anvil/model"
	"mutably/anvil/model/inflection"
	"mutably/anvil/parser"
	"mutably/anvil/parser/verb"
	"mutably/anvil/view"
	"os"
	"runtime"
)

// Run makes the application perform the requested operation.
// This value depends on the first argument passed to the program.
var Run func()

// init uses command-line flags to determine the value of Run().
func init() {
	flags := &AppFlags{}

	flag.IntVar(&flags.PageLimit, "limit", -1,
		"Limit import to processing N pages")
	flag.BoolVar(&flags.BeVerbose, "v", false,
		"Enable verbose logging")
	flag.StringVar(&flags.DBName, "d",
		"", "The database name")
	flag.StringVar(&flags.DBUser, "u", "",
		"The database user")
	flag.StringVar(&flags.DBPassword, "p", "",
		"The database user's password")
	flag.StringVar(&flags.DBHost, "host", "localhost",
		"The hostname of the database")
	flag.UintVar(&flags.DBPort, "port", 5432,
		"The database port")

	if len(os.Args) == 1 {
		Run = ShowHelp
	} else {
		Run = flags.GetIntent()
	}
}

func main() {
	Run()
}

// ShowHelp displays possible commands.
func ShowHelp() {
	fmt.Println(
		`Usage: anvil <command>

Commands:
* import
    - Imports an XML archive
* view
    - Views a specific page of an XML archive
* help
    - Displays information about command flags
	`)
}

// Import processes the contents of an archive.
func Import(args *AppFlags) {
	if flag.NArg() != 1 || args.MissingDBCredentials() {
		fmt.Println("Usage: anvil import -d [-h] [-port] -u -p <file>")
		os.Exit(1)
	}

	archive, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	psqlDB, err := model.NewPsqlDB(model.KeyRing{
		DatabaseName: args.DBName,
		Host:         args.DBHost,
		Port:         args.DBPort,
		User:         args.DBUser,
		Password:     args.DBPassword,
	})
	if err != nil {
		log.Fatal(err)
	}

	conjugators := inflection.NewConjugators()
	conjugators.Add(&inflection.Dutch{})

	vparser, err := verb.NewVerbParser(psqlDB, runtime.GOMAXPROCS(0),
		args.PageLimit, conjugators)
	if err != nil {
		log.Fatal(err)
	}

	parser.ProcessPages(archive, vparser)
	vparser.Wait()
}

// View displays content from the archive.
func View() {
	if flag.NArg() != 2 {
		fmt.Println("Usage: anvil view <file> <page title>")
	} else {
		view.Search(flag.Args()[0], flag.Args()[1])
	}
}
