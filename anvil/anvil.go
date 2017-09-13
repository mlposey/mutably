package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
// init() will assign a value to this var. That value depends
// on the first program argument.
var Run func()

// Parse CLI flags to determine what Run() should do.
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

	flag.Parse()

	if len(flag.Args()) == 0 {
		Run = func() {
			fmt.Println("Run anvil -h for usage details")
		}
	} else {
		if !flags.BeVerbose {
			log.SetOutput(ioutil.Discard)
		}

		// Determine the main purpose.
		switch flag.Args()[0] {
		case "import":
			Run = func() { Import(flags) }

		case "view":
			Run = View
		}
	}
}

func main() {
	Run()
}

// Import processes the contents of an archive.
func Import(args *AppFlags) {
	file, psqlDB := parseImportFlags(
		&args.DBName,
		&args.DBHost,
		&args.DBUser,
		&args.DBPassword,
		&args.DBPort,
	)

	conjugators := inflection.NewConjugators()
	conjugators.Add(&inflection.Dutch{})

	vparser, err := verb.NewVerbParser(psqlDB, runtime.GOMAXPROCS(0),
		args.PageLimit, conjugators)
	if err != nil {
		log.Fatal(err)
	}

	parser.ProcessPages(file, vparser)
	vparser.Wait()

}

// View displays content from the archive.
func View() {
	if len(flag.Args()) != 2 {
		fmt.Println("Usage: anvil view [file] [page-title]")
	} else {
		view.Search(flag.Args()[0], flag.Args()[1])
	}

}

func parseImportFlags(dbName, dbHost, dbUser, dbPwd *string,
	dbPort *uint) (*os.File, *model.PsqlDB) {
	if flag.NFlag() < 4 || *dbName == "" || *dbUser == "" || *dbPwd == "" ||
		len(flag.Args()) == 0 {
		fmt.Println("Usage: anvil import -d [-h] [-port] -u -p pages-file")
		os.Exit(1)
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err.Error())
	}

	key := model.KeyRing{
		DatabaseName: *dbName,
		Host:         *dbHost,
		Port:         *dbPort,
		User:         *dbUser,
		Password:     *dbPwd,
	}

	psqlDB, err := model.NewPsqlDB(key)
	if err != nil {
		log.Fatal(err.Error())
	}
	return file, psqlDB
}
