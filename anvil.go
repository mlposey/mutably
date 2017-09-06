package main

import (
	"anvil/model"
	"anvil/parse"
	"anvil/parse/verb"
	"anvil/view"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	profileFlag := flag.Bool("profile", false,
		`Enable program profiling.
		This uses Go's pprof library to profile program execution
		and output results to profile.txt. To inspect the output, call
		'go tool pprof profile.txt'.

		Executing 'web' in the profiler will generate a callgraph
		in your browser window. This feature requires graphviz, which
		can be installed on Ubuntu by running 'sudo apt install graphviz'.`)

	importFlag := flag.Bool("import", false,
		`Import a .xml pages dump into a PostgreSQL database.
		Takes the form -import -d [-host] [-port] -u -p pages-file
		The database (-d), user (-u), password (-p) and pages file are
		required flags. Host and port are optional and will default to
		'localhost' and 5432, respectively. The pages file should be
		pages dump without metadata, since import processes the raw
		xml file--not the compressed version. Metapages is far too
		large when uncompressed.

	Example:
		anvil -import -d=mutablydb -u=mutably -p=aPass latest-pages.xml`,
	)

	limitFlag := flag.Int("limit", -1,
		"Limit import to processing N pages")
	verboseFlag := flag.Bool("v", false,
		"Enable verbose logging")
	dbNameFlag := flag.String("d",
		"", "The database name")
	dbUserFlag := flag.String("u", "",
		"The database user")
	dbPwdFlag := flag.String("p", "",
		"The database user's password")
	dbHostFlag := flag.String("host", "localhost",
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

	if !*verboseFlag {
		log.SetOutput(ioutil.Discard)
	}

	if *profileFlag {
		profilingResults, _ := os.Create("profile.txt")
		pprof.StartCPUProfile(profilingResults)
		defer pprof.StopCPUProfile()
	}

	if *importFlag {
		file, psqlDB := parseImportFlags(dbNameFlag, dbHostFlag,
			dbUserFlag, dbPwdFlag, dbPortFlag)

		consumer, err := verb.NewVerbConsumer(psqlDB, runtime.GOMAXPROCS(0),
			*limitFlag)
		if err != nil {
			log.Fatal(err)
		}

		parse.ProcessPages(file, consumer)
		consumer.Wait()

	} else if *viewPageFlag {
		if len(flag.Args()) != 2 {
			fmt.Println("Usage: anvil -view [file] [page-title]")
		} else {
			view.Search(flag.Args()[0], flag.Args()[1])
		}
	}
}

func parseImportFlags(dbName, dbHost, dbUser, dbPwd *string,
	dbPort *uint) (*os.File, *model.PsqlDB) {
	if flag.NFlag() < 4 || *dbName == "" || *dbUser == "" || *dbPwd == "" ||
		len(flag.Args()) == 0 {
		fmt.Println("Usage: anvil -import -d [-h] [-port] -u -p pages-file")
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
