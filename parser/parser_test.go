package parser_test

import (
	"anvil/parser"
	"strings"
	"testing"
)

type testConsumer struct {
	Pages []parser.Page
}

func (c *testConsumer) Consume(page parser.Page) (bool, error) {
	c.Pages = append(c.Pages, page)
	return true, nil
}

// TODO: Generate a variable amount of pages.
// It would be nice if we could specify the page count and text content. Sha
// hashes, timestamps, id's, etc. would be generated uniquely for each page.

// Makes a sample page dump with 2 pages, returning the contents and page count.
func makePageDump() (string, int) {
	return `<mediawiki xmlns="http://www.mediawiki.org/xml/export-0.10/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.mediawiki.org/xml/export-0.10/ http://www.mediawiki.org/xml/export-0.10.xsd" version="0.10" xml:lang="en">
  <siteinfo>
    <sitename>Wiktionary</sitename>
    <dbname>enwiktionary</dbname>
    <base>https://en.wiktionary.org/wiki/Wiktionary:Main_Page</base>
    <generator>MediaWiki 1.30.0-wmf.7</generator>
    <case>case-sensitive</case>
    <namespaces>
      <namespace key="-2" case="case-sensitive">Media</namespace>
      <namespace key="-1" case="first-letter">Special</namespace>
      <namespace key="0" case="case-sensitive" />
      <namespace key="1" case="case-sensitive">Talk</namespace>
      <namespace key="2" case="first-letter">User</namespace>
      <namespace key="3" case="first-letter">User talk</namespace>
      <namespace key="4" case="case-sensitive">Wiktionary</namespace>
      <namespace key="5" case="case-sensitive">Wiktionary talk</namespace>
      <namespace key="6" case="case-sensitive">File</namespace>
      <namespace key="7" case="case-sensitive">File talk</namespace>
      <namespace key="8" case="first-letter">MediaWiki</namespace>
      <namespace key="9" case="first-letter">MediaWiki talk</namespace>
      <namespace key="10" case="case-sensitive">Template</namespace>
      <namespace key="11" case="case-sensitive">Template talk</namespace>
      <namespace key="12" case="case-sensitive">Help</namespace>
      <namespace key="13" case="case-sensitive">Help talk</namespace>
      <namespace key="14" case="case-sensitive">Category</namespace>
      <namespace key="15" case="case-sensitive">Category talk</namespace>
      <namespace key="90" case="case-sensitive">Thread</namespace>
      <namespace key="91" case="case-sensitive">Thread talk</namespace>
      <namespace key="92" case="case-sensitive">Summary</namespace>
      <namespace key="93" case="case-sensitive">Summary talk</namespace>
      <namespace key="100" case="case-sensitive">Appendix</namespace>
      <namespace key="101" case="case-sensitive">Appendix talk</namespace>
      <namespace key="102" case="case-sensitive">Concordance</namespace>
      <namespace key="103" case="case-sensitive">Concordance talk</namespace>
      <namespace key="104" case="case-sensitive">Index</namespace>
      <namespace key="105" case="case-sensitive">Index talk</namespace>
      <namespace key="106" case="case-sensitive">Rhymes</namespace>
      <namespace key="107" case="case-sensitive">Rhymes talk</namespace>
      <namespace key="108" case="case-sensitive">Transwiki</namespace>
      <namespace key="109" case="case-sensitive">Transwiki talk</namespace>
      <namespace key="110" case="case-sensitive">Wikisaurus</namespace>
      <namespace key="111" case="case-sensitive">Wikisaurus talk</namespace>
      <namespace key="114" case="case-sensitive">Citations</namespace>
      <namespace key="115" case="case-sensitive">Citations talk</namespace>
      <namespace key="116" case="case-sensitive">Sign gloss</namespace>
      <namespace key="117" case="case-sensitive">Sign gloss talk</namespace>
      <namespace key="118" case="case-sensitive">Reconstruction</namespace>
      <namespace key="119" case="case-sensitive">Reconstruction talk</namespace>
      <namespace key="828" case="case-sensitive">Module</namespace>
      <namespace key="829" case="case-sensitive">Module talk</namespace>
      <namespace key="2300" case="case-sensitive">Gadget</namespace>
      <namespace key="2301" case="case-sensitive">Gadget talk</namespace>
      <namespace key="2302" case="case-sensitive">Gadget definition</namespace>
      <namespace key="2303" case="case-sensitive">Gadget definition talk</namespace>
      <namespace key="2600" case="first-letter">Topic</namespace>
    </namespaces>
  </siteinfo>
  <page>
        <title>page1</title>
        <ns>0</ns>
        <id>1</id>
        <revision>
                <id>1</id>
                <timestamp>2017-07-01T09:18:05Z</timestamp>
                <contributor>
                        <username>tester</username>
                        <id>1</id>
                </contributor>
                <comment>A revision comment</comment>
                <model>wikitext</model>
                <format>text/x-wiki</format>
                <text xml:space="preserve">Sample text
                </text>
                <sha1>pe4p0zjy9806ye8ji3qah0p4pvcrdjt</sha1>
        </revision>
  </page>
  <page>
        <title>page2</title>
        <ns>1</ns>
        <id>2</id>
        <revision>
                <id>2</id>
                <timestamp>2017-07-01T09:17:45Z</timestamp>
                <contributor>
                        <username>tester</username>
                        <id>2</id>
                </contributor>
                <comment>A revision comment</comment>
                <model>wikitext</model>
                <format>text/x-wiki</format>
                <text xml:space="preserve">Sample text
                </text>
                <sha1>phm76it2loo31avsd7b9wykwqbb23wu</sha1>
        </revision>
  </page>
</mediawiki>`, 2
}

// ImportPages should read valid page dumps into their corresponding structs.
func TestImportPages(t *testing.T) {
	consumer := &testConsumer{}
	dump, pageCount := makePageDump()
	reader := strings.NewReader(dump)

	if err := parser.ProcessPages(reader, consumer); err != nil {
		t.Error(err.Error())
	}

	if len(consumer.Pages) != pageCount {
		t.Error("Expected", pageCount, "pages. Found", len(consumer.Pages))
	}

	const textContent = "Sample text"
	trim := func(s string) string { return strings.Trim(s, "\n ") }

	for _, page := range consumer.Pages {
		if trim(page.Revision.Text) != trim(textContent) {
			t.Error("Expected", textContent, ". Found",
				page.Revision.Text)
		}
	}
}
