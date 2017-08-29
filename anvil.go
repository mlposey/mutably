package main

import (
	"flag"
	"fmt"
	"anvil/view"
	"anvil/db"
	"anvil/parse"
	"os"
)

func main() {
	importFlag := flag.Bool("import", false,
		"Import a .xml pages dump.")
	dbNameFlag := flag.String("d",
		"", "The database name")
	dbUserFlag := flag.String("u", "",
		"The database user")
	dbPwdFlag := flag.String("p", "",
		"The database user's password")
	dbHostFlag := flag.String("h", "localhost",
		"The hostname of the database")
	dbPortFlag := flag.Uint("port", 5432,
		"The database port")

	viewPageFlag := flag.Bool("view", false,
		"View a page by its title.")

	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Run 'anvil -h' for usage details.")
		return
	}

	if *importFlag {
		file, credentials := parseImportFlags(dbNameFlag, dbHostFlag,
			dbUserFlag, dbPwdFlag, dbPortFlag)
		parse.ProcessPages(file, &parse.VerbConsumer{credentials})

	} else if *viewPageFlag {
		if len(flag.Args()) != 2 {
			fmt.Println("Usage: anvil -view [file] [page-title]")
		} else {
			view.Search(flag.Args()[0], flag.Args()[1])
		}
	}
}

func parseImportFlags(dbName, dbHost, dbUser, dbPwd *string, dbPort *uint) (*os.File, db.KeyRing) {
	if flag.NFlag() < 4 || *dbName == "" || *dbUser == "" || *dbPwd == "" ||
		len(flag.Args()) == 0 {
		fmt.Println("Usage: anvil -import -d [-h] [-port] -u -p pages-file")
		os.Exit(1)
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Println("Failed to open file", flag.Arg(0))
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return file, db.KeyRing{
		Database: *dbName,
		Host: 	  *dbHost,
		Port: 	  *dbPort,
		User: 	  *dbUser,
		Password: *dbPwd,
	}
}
