package parser

import "encoding/xml"

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
