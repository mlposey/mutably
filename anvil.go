package main

import (
	"flag"
	"fmt"
	"anvil/view"
)

func main() {
	importFlag := flag.Bool("import", false, "Import a .xml pages dump.")
	viewPageFlag := flag.Bool("view", false, "View a page by its title.")
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Run 'anvil -h' for usage details.")
		return
	}

	if *importFlag {
		fmt.Println("Import stub")
	} else if *viewPageFlag {
		if len(flag.Args()) != 2 {
			fmt.Println("Usage: anvil -view [file] [page-title]")
		} else {
			view.Search(flag.Args()[0], flag.Args()[1])
		}
	}
}
