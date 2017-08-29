package view

import (
	"anvil/parse"
	"os"
	"fmt"
	"github.com/fatih/color"
)

type PageViewer struct {
	TargetPageTitle string
}
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
		fmt.Println("Could not open file", filePath)
	}

	color.Green("Starting search...")
	err = parse.ProcessPages(file, &PageViewer{title})
	if err != nil {
		fmt.Println(err.Error())
	}
	color.Green("Done.")
}
