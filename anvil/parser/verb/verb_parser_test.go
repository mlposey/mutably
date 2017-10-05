package verb_test

import (
	"mutably/anvil/model"
	"mutably/anvil/model/inflection"
	"mutably/anvil/parser"
	"mutably/anvil/parser/verb"
	"testing"
)

func TestParse(t *testing.T) {
	vparser, mdb := makeMockParser(t)

	cont, err := vparser.Parse(mockPage.Page)
	vparser.Wait()

	if cont == false {
		t.Error("Continue signal should be true")
	}
	if err != nil {
		t.Error(err.Error())
	}

	// We subtract one because English defines a verb twice and detecting
	// multiple verb definitions is on halt for now.
	if len(mdb.verbs) != mockPage.TemplateCount-1 {
		t.Error("Expected", mockPage.VerbCount, "verbs, found", len(mdb.verbs))
	}
}

// VerbParser should not process a section of a page if it is for a language
// that is undefined.
func TestVerbParser_NewLanguage(t *testing.T) {
	emptyConjugators := make(map[string]inflection.Conjugator)
	mdb := newMockDB()
	parser, e := verb.NewVerbParser(mdb, 2, -1, emptyConjugators)
	if e != nil {
		t.Error(e.Error())
	}

	parser.Parse(mockPage.Page)
	parser.Wait()

	if len(mdb.words) != 0 || len(mdb.verbs) != 0 {
		t.Error("Unidentified languages should not be processed.")
	}
}

func makeMockParser(t *testing.T) (*verb.VerbParser, *mockDB) {
	t.Helper()

	// Make a conjugator for each language in the mock page.
	conjugators := make(map[string]inflection.Conjugator)
	for _, language := range mockPage.Languages {
		conjugators[language.String()] = &mockConjugator{language: language}
	}

	db := newMockDB()

	vp, e := verb.NewVerbParser(db, 2, -1, conjugators)
	if e != nil {
		t.Error(e.Error())
	}
	return vp, db
}

type mockConjugator struct {
	language *model.Language
	database model.Database
}

func (m *mockConjugator) GetLanguage() *model.Language {
	return m.language
}
func (m *mockConjugator) SetDatabase(db model.Database) error {
	m.database = db
	return nil
}
func (m *mockConjugator) Conjugate(verb, template string) error {
	err := m.database.InsertVerbForm(&model.VerbForm{Word: verb})
	return err
}

type mockDB struct {
	languages []*model.Language
	words     []string
	verbs     []*model.VerbForm
}

func newMockDB() *mockDB {
	return &mockDB{languages: mockPage.Languages}
}

func (m *mockDB) InsertLanguage(language *model.Language) error {
	m.languages = append(m.languages, language)
	return nil
}
func (m *mockDB) InsertWord(word string) (wordId int) {
	m.words = append(m.words, word)
	return len(m.words) - 1
}
func (m *mockDB) InsertVerbForm(verb *model.VerbForm) error {
	m.verbs = append(m.verbs, verb)
	return nil
}

// A mock page setup. Do not change the contents of Page! Many tests in this
// test package rely on it being the way it is. If you need a different
// structure, make your own.
var mockPage = struct {
	Languages     []*model.Language // The languages in the page
	VerbCount     int               // The number of languages where the word is a verb
	TemplateCount int               // The total number of verb templates
	Page          parser.Page
}{
	Languages: []*model.Language{
		model.NewLanguage("English"),
		model.NewLanguage("spanish"),
		model.NewLanguage("Finnish"),
		model.NewLanguage("french"),
	},
	VerbCount:     4,
	TemplateCount: 14,
	Page: parser.Page{
		Title: "lie",
		Revision: parser.Revision{
			Text: `
		{{also|LIE|lié|líe|liè|liē|liě|li'e}}
{{TOC limit|3}}
==English==

====Verb====
{{en-verb|lies|lying|lay|lain}}

# {{senseid|en|to rest}} {{lb|en|intransitive}} To [[rest#Verb|rest]] in a [[horizontal]] [[position]] on a [[surface]].

# {{lb|en|legal}} To be [[sustainable]]; to be capable of being [[maintain]]ed.
#* Ch. J. Parsons
#*: An appeal '''lies''' in this case.

====Noun====
{{en-noun}}

# {{lb|en|golf}} The [[terrain]] and [[condition]]s surrounding the [[ball]] before it is [[strike#Verb|struck]].
# {{lb|en|medicine}} The position of a [[fetus]] in the [[womb]].

=====Translations=====
{{trans-top|golf term}}
* Catalan: {{t+|ca|situació|f}}
* Dutch: {{t+|nl|ligging|f}}, {{t|nl|terreinligging|f}}
{{trans-mid}}

====Verb====
{{en-verb|lies|lying|lied}}

# {{senseid|en|false}} {{lb|en|intransitive}} To give [[false]] [[information]] [[intentional]]ly.
#: {{ux|en|Hips don't '''lie'''.}}

=====Synonyms=====
* {{l|en|prevaricate}}

====Noun====
{{en-noun}}

# An [[intentionally]] [[false]] [[statement]]; an [[intentional]] [[falsehood]].
#*: Wishing this '''lie''' of life was o'er.
#* ''The cake is a '''lie'''.'' - [[w:Portal (video game)|Portal]]

=====Synonyms=====
{{top2}}
* {{l|en|alternative fact}}
* {{l|en|falsehood}}

----

==Finnish==

===Verb===
{{head|fi|verb form}}

# {{lb|fi|nonstandard}} {{fi-form of|olla|pr=third-person|pl=singular|mood=potential|tense=present}}
#: ''Se on missä '''lie'''.''
#:: ''It's somewhere.'' / ''I wonder where it is.''
#: ''Tai mitä '''lie''' ovatkaan''
#:: ''Or whatever they are.''

===Anagrams===
* {{l|fi|eli}}, {{l|fi|lei}}

----

==French==

===Verb===
{{fr-verb-form}}

# {{inflection of|lier||1|s|pres|indc|lang=fr}}
# {{inflection of|lier||3|s|pres|indc|lang=fr}}
# {{inflection of|lier||1|s|pres|subj|lang=fr}}
# {{inflection of|lier||3|s|pres|subj|lang=fr}}
# {{inflection of|lier||2|s|impr|lang=fr}}

===Anagrams===
* {{l|fr|île}}

===Further reading===
* {{R:TLFi}}

----

==Spanish==

===Verb===
{{head|es|verb form}}

# {{es-verb form of|ending=ar|mood=subjunctive|tense=present|pers=1|number=singular|liar}}
# {{es-verb form of|ending=ar|mood=subjunctive|tense=present|pers=3|number=singular|liar}}
	`}},
}
