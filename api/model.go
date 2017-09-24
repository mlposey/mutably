package main

// NewErrorResponse creates an error message in the APIs standard format.
func NewErrorResponse(msg string) map[string]string {
	return map[string]string{
		"error": msg,
	}
}

// Language describes a natural language that exists in the database.
type Language struct {
	// The row id of the language
	Id int `json:"id"`
	// The name of the language (e.g., 'English')
	Name string `json:"name"`
	// The language tag (e.g., 'en')
	Tag string `json:"tag"`
}

// GetLanguage returns from the database the language identified by id.
func (db *PsqlDB) GetLanguage(id int) (*Language, error) {
	language := &Language{Id: id}

	err := db.QueryRow(`
		SELECT language, tag FROM languages
		WHERE id = $1`,
		language.Id,
	).Scan(&language.Name, &language.Tag)

	if err != nil {
		return nil, err
	}
	return language, nil
}

// GetLanguages returns a slice of all languages in the database.
func (db *PsqlDB) GetLanguages() ([]*Language, error) {
	rows, err := db.Query(`SELECT * FROM languages`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var languages []*Language
	for rows.Next() {
		lang := &Language{}

		err = rows.Scan(&lang.Id, &lang.Name, &lang.Tag)
		if err != nil {
			return nil, err
		}
		languages = append(languages, lang)
	}

	return languages, nil
}

// Word models a word from a specific language.
type Word struct {
	// The id of the row containing the word
	Id int `json:"id"`
	// The word itself
	Text string `json:"text"`
	// The id of the word's language
	LanguageId int `json:"language"`
}

// GetWords returns a slice of all words in the database.
func (db *PsqlDB) GetWords() ([]*Word, error) {
	rows, err := db.Query(`
		SELECT words.id, words.word, languages.id
		FROM verbs
		JOIN words     on verbs.word_id = words.id
		JOIN languages on verbs.lang_id = languages.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*Word
	for rows.Next() {
		word := &Word{}
		err = rows.Scan(&word.Id, &word.Text, &word.LanguageId)
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}

// GetWord returns from the database a word identified by id.
func (db *PsqlDB) GetWord(id int) (*Word, error) {
	word := &Word{}
	err := db.QueryRow(`
		SELECT words.id, words.word, languages.id
		FROM verbs
		JOIN words     on verbs.word_id = $1
		JOIN languages on verbs.lang_id = languages.id`,
		id,
	).Scan(&word.Id, &word.Text, &word.LanguageId)

	if err != nil {
		return nil, err
	}
	return word, nil
}
