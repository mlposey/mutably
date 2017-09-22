package parser

import (
	"encoding/xml"
	"io"
	"strings"
)

// A Parser processes the contents of a Page.
type Parser interface {
	Parse(Page) (bool, error)
}

// ProcessPages sends each page of pagesFile to a parser.
//
// See the Page documentation for details on expected structure for pagesFile.
func ProcessPages(pagesFile io.Reader, parser Parser) error {
	decoder := xml.NewDecoder(pagesFile)
	for t, e := decoder.Token(); t != nil; t, e = decoder.Token() {
		if e != nil {
			return e
		}

		switch elementType := t.(type) {
		case xml.StartElement:
			{
				if elementType.Name.Local == "page" {
					var page Page
					decoder.DecodeElement(&page, &elementType)

					if isSpecialPage(&page.Title) {
						continue
					}

					cont, err := parser.Parse(page)
					if err != nil {
						return err
					}
					if cont == false {
						return nil
					}
				}
			}
		}
	}
	return nil
}

func isSpecialPage(pageTitle *string) bool {
	colonPos := strings.IndexRune(*pageTitle, ':')
	if colonPos == -1 {
		return false
	}

	switch (*pageTitle)[0:colonPos] {
	case "User", "Template", "Module", "Category", "Thread", "Talk":
		return true
	default:
		return false
	}
}
