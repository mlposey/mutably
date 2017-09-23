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

		err := rows.Scan(&lang.Id, &lang.Name, &lang.Tag)
		if err != nil {
			return nil, err
		}
		languages = append(languages, lang)
	}

	return languages, nil
}
