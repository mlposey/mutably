package parse

import (
	"encoding/xml"
	"io"
)

// Page defines the XML structure of a wiktionary 'pages' export.
// Most pages contain information about a word for various languages. A few
// others describe templates that are used by those languages.
//
// Page relies on a reduced form of the schema defined at
// https://www.mediawiki.org/xml/export-0.8.xsd. Certain exports contain full
// revision history. Here, we are only concerned with the most recent version
// of a page.
type Page struct {
	XMLName xml.Name `xml:"page"`

	Title    string   `xml:"title"`
	Revision Revision `xml:"revision"`
}

// Revision defines the XML structure for a version of a page.
// Certain exports contain multiple revisions, which are copies of a page at
// each point of its history of modification. That extra data is unneeded.
// Because we only need one revision, this struct is essentially the contents
// of a Page.
type Revision struct {
	XMLName xml.Name `xml:"revision"`

	// The contents of the revision (i.e., body of the webpage)
	Text string `xml:"text"`
}

// A PageConsumer processes Page objects.
type PageConsumer interface {
	Consume(page Page) (bool, error)
}

// ProcessPages passes each page element of the pagesFile to consumer.
//
// consumer.Consume should return true to consume more pages or false
// to stop.
//
// See the Page documentation for details on expected structure for pagesFile.
func ProcessPages(pagesFile io.Reader, consumer PageConsumer) error {
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

					cont, err := consumer.Consume(page)
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
