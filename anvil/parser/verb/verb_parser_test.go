package verb_test

import (
	"mutably/anvil/model"
	"mutably/anvil/parser"
	"mutably/anvil/parser/verb"
	"testing"
)

func TestParse(t *testing.T) {
	mdb := NewMockDB()
	vparser, e := verb.NewVerbParser(mdb, 2, -1)
	if e != nil {
		t.Error(e.Error())
	}

	cont, err := vparser.Parse(mockPage.Page)
	vparser.Wait()

	if cont == false {
		t.Error("Continue signal should be true")
	}
	if err != nil {
		t.Error(err.Error())
	}

	if len(mdb.verbs) != mockPage.VerbCount {
		t.Error("Expected", mockPage.VerbCount, "verbs, found", len(mdb.verbs))
	}

	if mdb.TemplateCount() != mockPage.TemplateCount {
		t.Error("Expected", mockPage.TemplateCount, "templates, found",
			mdb.TemplateCount())
	}
}

// VerbParser should not process a section of a page if it is for a language
// that is undefined.
func TestVerbParser_NewLanguage(t *testing.T) {
	mdb := NewMockDB()
	mdb.languages = []model.Language{}

	parser, e := verb.NewVerbParser(mdb, 2, -1)
	if e != nil {
		t.Error(e.Error())
	}

	parser.Parse(mockPage.Page)
	parser.Wait()

	if len(mdb.words) != 0 || len(mdb.verbs) != 0 || len(mdb.templates) != 0 {
		t.Error("Unidentified languages should not be processed.")
	}
}

// VerbParser should store both versions of a verb that has multiple meanings.
// E.g., 'lie', which has two meanings in English, will have two verb
// sections in the archive. These should produce two different verb entries
// in the database.
func TestVerbParser_MultipleMeanings(t *testing.T) {
	mdb := NewMockDB()
	parser, e := verb.NewVerbParser(mdb, 2, -1)
	if e != nil {
		t.Error(e.Error())
	}
	parser.Parse(mockPage.Page)
	parser.Wait()

	lieIndex := -1
	for i := range mdb.words {
		if mdb.words[i] == "lie" {
			lieIndex = i
			break
		}
	}
	if lieIndex == -1 {
		t.Error("VerbParser is not inserting words of known languages.")
	}

	if len(mdb.verbs[lieIndex]) != 2 {
		t.Error("VerbParser is not adding duplicate verbs with different meanings.")
	}
}

type mockDB struct {
	languages []model.Language
	words     []string
	verbs     map[int][]int                // languageId -> {wordId...}
	templates map[model.VerbTemplate][]int // template -> {verbId...}
}

func NewMockDB() *mockDB {
	return &mockDB{
		languages: []model.Language{"english", "spanish", "finnish", "french"},
		verbs:     make(map[int][]int),
		templates: make(map[model.VerbTemplate][]int),
	}
}

func (mdb *mockDB) LanguageExists(lang model.Language) (bool, int) {
	for i, language := range mdb.languages {
		if language == lang {
			return true, i
		}
	}
	return false, -1
}

func (mdb *mockDB) InsertWord(word string) int {
	mdb.words = append(mdb.words, word)
	return len(mdb.words) - 1
}

func (mdb *mockDB) InsertVerb(wordId, languageId int) (int, error) {
	mdb.verbs[languageId] = append(mdb.verbs[languageId], wordId)
	return len(mdb.verbs[languageId]) - 1, nil
}

func (mdb *mockDB) InsertTemplate(template model.VerbTemplate, verbId int) error {
	// This doesn't really make a whole ton of sense, but it doesn't need
	// to for current testing. A more accurate representation of template
	// storage may be needed if the complexity of VerbParser tests increases.
	mdb.templates[template] = append(mdb.templates[template], verbId)
	return nil
}

func (mdb *mockDB) TemplateCount() (count int) {
	for _, v := range mdb.templates {
		count += len(v)
	}
	return
}

// A mock page setup. Do not change the contents of Page! Many tests in this
// test package rely on it being the way it is. If you need a different
// structure, make your own.
var mockPage = struct {
	LanguageCount int // The number of languages on the page
	VerbCount     int // The number of languages where the word is a verb
	TemplateCount int // The total number of verb templates
	Page          parser.Page
}{
	LanguageCount: 4,
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
