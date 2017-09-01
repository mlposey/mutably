package parse_test

import (
	"testing"
	"anvil/parse"
	"fmt"
)

// *VerbConsumer.Consume should split a page into sections, where
// each section is the word's context for a specific language.
func TestSectionDivision(t *testing.T) {
	counter := parse.NewVerbConsumer(nil)
	counter.Consume(mockPage.Page)

	if counter.LanguageCount != mockPage.LanguageCount {
		t.Error("Expected", mockPage.LanguageCount, "got",
			counter.LanguageCount)
	}
}

// *VerbConsumer.Consume should identify which languages define
// the word as a verb.
func TestCountVerbs(t *testing.T) {
	counter := parse.NewVerbConsumer(nil)
	counter.Consume(mockPage.Page)

	if counter.VerbCount != mockPage.VerbCount {
		t.Error("Expected", mockPage.VerbCount, "got",
			counter.VerbCount)
	}
}

// GetTemplates should make a *Verb out of each template a language defines.
func TestGetTemplates(t *testing.T) {
	word := "blepsh"
	language := "Blosh"

	templates := []string{
		"{{test-template|test}}",
		"{{super-test-template|test}}",
	}
	content := fmt.Sprintf(
		`
		test test test
		test

		====Verb====
		%s

		%s

		test
		test test

		====Verb====
		%s

		test
		`, templates[0], templates[0], templates[1])

	consumer := parse.NewVerbConsumer(nil)
	consumer.CurrentSection = content
	verbs := consumer.GetTemplates(&word, &language)

	if len(verbs) != len(templates) {
		t.Error("Expected", len(templates), "templates, found", len(verbs))
	}

	if
		verbs[0].Template == verbs[1].Template ||
		(verbs[0].Template != templates[0] && verbs[0].Template != templates[1]) ||
		(verbs[1].Template != templates[0] && verbs[1].Template != templates[1]) {
		t.Error("Failed to read templates")
	}
}

var mockPage = struct {
	LanguageCount int // The number of languages on the page
	VerbCount     int // The number of languages where the word is a verb
	Page          parse.Page
}{
	LanguageCount: 7,
	VerbCount: 4,
	Page: parse.Page{
		Title: "lie",
		Revision: parse.Revision{
			Text:
			`
		{{also|LIE|lié|líe|liè|liē|liě|li'e}}
{{TOC limit|3}}
==English==

===Pronunciation===
* {{IPA|/laɪ̯/|lang=en}}
* {{audio|en-us-lie.ogg|Audio (GA)|lang=en}}
* {{rhymes|aɪ|lang=en}}
* {{homophones|lang=en|lye|lai}}

===Etymology 1===
{{PIE root|en|legʰ}}
From {{inh|en|enm|lien}}, {{m|enm|liggen}}, from {{inh|en|ang|licgan}}, from {{inh|en|gem-pro|*ligjaną}}, from {{der|en|ine-pro|*legʰ-}}. Cognate with {{cog|fy|lizze}}, {{cog|nl|liggen}}, {{cog|de|liegen}}, {{cog|da|-}} and {{cog|nb|ligge}}, {{cog|sv|ligga}}, {{cog|nn|liggja}}, {{cog|got|𐌻𐌹𐌲𐌰𐌽}}; and with {{cog|la|lectus||bed}}, {{cog|ga|luighe}}, {{cog|ru|лежа́ть}}, {{cog|sq|lagje||inhabited area, neighbourhood}}.

As a noun for {{m|en|position}}, the [[#Noun|noun]] has the same etymology above as the [[#Verb|verb]].

====Verb====
{{en-verb|lies|lying|lay|lain}}

# {{senseid|en|to rest}} {{lb|en|intransitive}} To [[rest#Verb|rest]] in a [[horizontal]] [[position]] on a [[surface]].

# {{lb|en|legal}} To be [[sustainable]]; to be capable of being [[maintain]]ed.
#* Ch. J. Parsons
#*: An appeal '''lies''' in this case.

=====Usage notes=====
The verb ''lie'' in this sense is sometimes used interchangeably with the verb {{m|en|lay}} in informal spoken settings. This can lead to nonstandard constructions which are sometimes objected to. Additionally, the past tense and past participle can both become {{m|en|laid}}, instead of {{m|en|lay}} and {{m|en|lain}} respectively, in less formal settings. These usages are common in speech but rarely found in edited writing or in more formal spoken situations.

=====Derived terms=====
{{der3|lang=en|[[a lie has no legs]]
|[[let sleeping dogs lie]]
|[[lie back]]
|[[lie by]]
|[[make one's bed and lie in it]]
|[[therein lies the rub]]
}}

=====Related terms=====
* [[lay#Etymology 1|lay]], a corresponding transitive version of this word
* {{l|en|lees}}
* {{l|en|lier}}

=====Translations=====
{{trans-top|be in horizontal position}}
* Afrikaans: {{t|af|lê}}
* Walloon: {{t|wa|esse metou|m}}
{{trans-bottom}}

====Noun====
{{en-noun}}

# {{lb|en|golf}} The [[terrain]] and [[condition]]s surrounding the [[ball]] before it is [[strike#Verb|struck]].
# {{lb|en|medicine}} The position of a [[fetus]] in the [[womb]].

=====Translations=====
{{trans-top|golf term}}
* Catalan: {{t+|ca|situació|f}}
* Dutch: {{t+|nl|ligging|f}}, {{t|nl|terreinligging|f}}
{{trans-mid}}
* Japanese: {{t+|ja|ライ|tr=rai}}
{{trans-bottom}}

{{trans-top|position of fetus}}
* Dutch: {{t+|nl|ligging|f}}
* German: {{t|de|Kindslage|f}}
{{trans-mid}}
* Norwegian: {{t+|no|leie|n}}
{{trans-bottom}}

===Etymology 2===
{{PIE root|en|lewgʰ}}
From {{inh|en|enm|lien||to lie, tell a falsehood}}, from {{inh|en|ang|lēogan||to lie}}, from {{inh|en|gem-pro|*leuganą||to lie}}, from {{der|en|ine-pro|*lewgʰ-||to lie, swear, bemoan}}. Cognate with {{cog|fy|lige||to lie}}, {{cog|nds|legen}}, {{m|nds|lögen||to lie}}, {{cog|nl|liegen||to lie}}, {{cog|de|lügen||to lie}}, {{cog|no|ljuge}}/{{m|no|lyge||to lie}}, {{cog|da|lyve||to lie}}, {{cog|sv|ljuga||to lie}}, and more distantly with {{cog|bg|лъжа||to lie}}, {{cog|ru|лгать||to lie}}, {{m|ru|ложь||falsehood}}.

====Verb====
{{en-verb|lies|lying|lied}}

# {{senseid|en|false}} {{lb|en|intransitive}} To give [[false]] [[information]] [[intentional]]ly.
#: {{ux|en|Hips don't '''lie'''.}}

=====Synonyms=====
* {{l|en|prevaricate}}

=====Derived terms=====
* {{l|en|liar}}
* {{l|en|lie through one's teeth}}

=====Translations=====
{{trans-top|tell an intentional untruth}}
* Abkhaz: {{t-needed|ab}}
* Afrikaans: {{t|af|lieg}}, {{t|af|jok}}
* Albanian: {{t+|sq|gënjej}}
* Walloon: {{t+|wa|minti}}, {{t+|wa|bourder}}
* West Frisian: {{t|fy|lige}}
* Yiddish: {{t|yi|לײַגן}}
{{trans-bottom}}

===Etymology 3===
From {{inh|en|enm|lie}}, from {{inh|en|ang|lyġe||lie, falsehood}}, from {{inh|en|gem-pro|*lugiz||lie, falsehood}}, from {{der|en|ine-pro|*leugh-||to tell lies, swear, complain}}, {{m|ine-pro|*lewgʰ-}}. Cognate with {{cog|osx|luggi||a lie}}, {{cog|goh|lugi|lugī}}, {{m|goh|lugin||a lie}} ({{cog|de|Lüge}}), {{cog|da|løgn||a lie}}, {{cog|bg|лъжа́||а lie}}.

====Noun====
{{en-noun}}

# An [[intentionally]] [[false]] [[statement]]; an [[intentional]] [[falsehood]].
#*: Wishing this '''lie''' of life was o'er.
#* ''The cake is a '''lie'''.'' - [[w:Portal (video game)|Portal]]

=====Synonyms=====
{{top2}}
* {{l|en|alternative fact}}
* {{l|en|falsehood}}
{{mid2}}
* {{l|en|fib}}
{{bottom}}
* See also [[Wikisaurus:lie]]

=====Antonyms=====
* {{l|en|truth}}

=====Derived terms=====
{{der3|lang=en|[[barefaced lie]]
|[[belie]]
|[[white lie]]
}}

=====Translations=====
{{trans-top|intentionally false statement}}
* Afrikaans: {{t|af|leuen}}
* Albanian: {{t+|sq|gënjeshtër|f}}
* Yiddish: {{t|yi|ליגן|m}}
{{trans-bottom}}

====Statistics====
* {{rank|turning|village|quickly|814|lie|supposed|original|provide}}

===Further reading===
* {{pedia}}

===Anagrams===
* {{anagrams|en|%ile|-ile|Eli|Ile|ile|lei}}

[[Category:English basic words]]
[[Category:English irregular verbs]]
[[Category:English terms with multiple etymologies]]
[[Category:English words following the I before E except after C rule]]

----

==Finnish==

===Verb===
{{head|fi|verb form}}

# {{lb|fi|nonstandard}} {{fi-form of|olla|pr=third-person|pl=singular|mood=potential|tense=present}}
#: ''Se on missä '''lie'''.''
#:: ''It's somewhere.'' / ''I wonder where it is.''
#: ''Tai mitä '''lie''' ovatkaan''
#:: ''Or whatever they are.''

====Usage notes====
* This form is chiefly used in direct and indirect questions.

====Synonyms====
* (''3rd-pers. sg. potent. pres. of olla; standard'') [[lienee]]

===Anagrams===
* {{l|fi|eli}}, {{l|fi|lei}}

----

==French==

===Etymology===
Probably from {{der|fr|xtg||*liga|silt, sediment}}, from {{der|fr|ine-pro|*legʰ-||to lie, to lay}}.

===Noun===
{{fr-noun|f}}

# [[lees]], [[dregs]] (of wine, of society)

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

==Mandarin==

===Romanization===
{{cmn-pinyin}}

# {{pinyin reading of|咧}}

{{head|cmn|pinyin}}

# {{nonstandard spelling of|lang=cmn|sc=Latn|liē}}
# {{nonstandard spelling of|lang=cmn|sc=Latn|lié}}
# {{nonstandard spelling of|lang=cmn|sc=Latn|liě}}
# {{nonstandard spelling of|lang=cmn|sc=Latn|liè}}

====Usage notes====
* {{cmn-toneless-note}}

----

==Old French==

===Etymology===
See English {{m|en|lees}}.

===Noun===
{{fro-noun|f}}

# [[dregs]]; mostly solid, undesirable leftovers of a drink

====Descendants====
* English: {{l|en|lees}}

----

==Spanish==

===Verb===
{{head|es|verb form}}

# {{es-verb form of|ending=ar|mood=subjunctive|tense=present|pers=1|number=singular|liar}}
# {{es-verb form of|ending=ar|mood=subjunctive|tense=present|pers=3|number=singular|liar}}

----

==Swedish==

===Etymology===
From {{etyl|gmq-osw|sv}} {{m|gmq-osw|līe}}, {{m|gmq-osw|lē}}, from {{etyl|non|sv}} {{m|non|lé}}, from {{etyl|gem-pro|sv}} {{m|gem-pro|*lewą}}, from {{etyl|ine-pro|sv}} {{m|ine-pro|*leu-||to cut}}.

===Pronunciation===
* {{IPA|/liːɛ/|lang=sv}}

===Noun===
{{sv-noun|c}}

# [[scythe]]; an instrument for mowing grass, grain, or the like.

====Declension====
{{sv-infl-noun-c-ar|2=li}}

====Related terms====
* {{l|sv|lieblad}}
* {{l|sv|lietag}}

===References===
* {{R:SAOL}}
		`},
	},
}
