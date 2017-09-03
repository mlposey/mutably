package view

import (
	"anvil/parse"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
)

// A PageViewer consumes Pages until one with a target title is found.
type PageViewer struct {
	TargetPageTitle string
}

// Consume checks page for a specific title. If the title is found, the
// page contents are printed to stdout and Consume returns false.
func (v *PageViewer) Consume(page parse.Page) (bool, error) {
	if page.Title == v.TargetPageTitle {
		fmt.Println(page)
		return false, nil
	}
	return true, nil
}

// Search looks for a page in filePath that contains title.
// If found, the page is printed to stdout.
func Search(filePath, title string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	color.Green("Starting search...")
	err = parse.ProcessPages(file, &PageViewer{title})
	if err != nil {
		fmt.Println(err.Error())
	}
	color.Green("Done.")
}
