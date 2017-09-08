package view

import (
	"anvil/parser"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
)

// PageViewer reads Pages until one with a target title is found.
type PageViewer struct {
	TargetPageTitle string
}

// Parse checks page for a specific title. If the title is found, the
// page contents are printed to stdout and Parse returns false.
func (v *PageViewer) Parse(page parser.Page) (bool, error) {
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
	err = parser.ProcessPages(file, &PageViewer{title})
	if err != nil {
		fmt.Println(err.Error())
	}
	color.Green("Done.")
}
