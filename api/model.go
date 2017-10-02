package main

import (
	"database/sql"
	"errors"
	"time"
)

// Language describes a natural language that exists in the database.
type Language struct {
	// The row id of the language
	Id int `json:"id"`

	// The name of the language (e.g., 'English')
	Name string `json:"name"`

	// The language tag (e.g., 'en')
	Tag sql.NullString `json:"tag"`
}

// GetLanguage returns from the database the language identified by id.
func (db *PsqlDB) GetLanguage(id int) (*Language, error) {
	language := &Language{Id: id}

	err := db.QueryRow(`
		SELECT name, tag FROM languages
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
		SELECT DISTINCT words.id, words.word, lang_id
		FROM verb_forms
		JOIN words
		  on verb_forms.word_id = words.id
		  or verb_forms.inf_id  = words.id
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
// TODO: This won't work for multiple languages.
func (db *PsqlDB) GetWord(id int) (*Word, error) {
	word := &Word{Id: id}
	err := db.QueryRow(`
		SELECT words.word, lang_id
		FROM verb_forms
		JOIN words on words.id = verb_forms.word_id
		WHERE word_id = $1`,
		id,
	).Scan(&word.Text, &word.LanguageId)

	if err != nil {
		return nil, err
	}
	return word, nil
}

// User models a user account as found in the database.
type User struct {
	// A UUID for the user
	Id string

	// The user's handle/name
	Name string

	// The id of the user's role (e.g., admin -> 1, user -> 2)
	RoleId int

	// The id of the language that the user prefers to work with
	TargetLanguageId sql.NullInt64

	// A timestamp dated when the user was created
	CreatedAt time.Time
}

// GetUsers returns a slice of all users in the database.
func (db *PsqlDB) GetUsers() ([]*User, error) {
	rows, err := db.Query(`
		SELECT id, role_id, name, target_language_id, created_at
		FROM users;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err = rows.Scan(&user.Id, &user.RoleId, &user.Name,
			&user.TargetLanguageId, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUser returns from the database a user identified by id.
func (db *PsqlDB) GetUser(id string) (*User, error) {
	user := &User{}
	err := db.QueryRow(`
		SELECT id, role_id, name, target_language_id, created_at
		FROM users
		WHERE id = $1`,
		id,
	).Scan(&user.Id, &user.RoleId, &user.Name, &user.TargetLanguageId,
		&user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser inserts a new user into the database and returns their id.
// Duplicate names are not allowed; attempting to insert one will result
// in a non-nil error.
func (db *PsqlDB) CreateUser(name, password string) (string, error) {
	var userId string
	err := db.QueryRow(`SELECT create_user($1, $2)`, name, password).Scan(&userId)
	return userId, err
}

// IsAdmin returns true if the user has administrator pivileges.
func (db *PsqlDB) IsAdmin(userId string) bool {
	var role string
	err := db.QueryRow(`
		SELECT role FROM roles
		WHERE id = (
			SELECT role_id FROM users
			WHERE id = $1
		)`, userId,
	).Scan(&role)
	if err != nil {
		return false
	}

	return role == "admin"
}

// Conjugation table stores the present and past tense forms of an infintive.
type ConjugationTable struct {
	Infinitive string
	Present    *TenseInflection
	Past       *TenseInflection
}

// TenseInflection stores the forms of a verb in a certain tense.
type TenseInflection struct {
	First  []string
	Second []string
	Third  []string
	Plural []string
}

// Person defines the grammatical person of a finite verb form.
// TODO: anvil has a similar definition. Try to share them.
type Person int

const (
	First  Person = 1 << 1
	Second Person = 1 << 2
	Third  Person = 1 << 3
)

// GetConjugationTable retrieves a tense inflection for word.
func (db *PsqlDB) GetConjugationTable(word string) (*ConjugationTable, error) {
	inf, infId, err := db.GetInfinitive(word)
	if err == sql.ErrNoRows {
		return nil, errors.New("word " + word + " does not exist")
	} else if err != nil {
		return nil, err
	}

	// We won't read person into a nullable type, so it is important that the
	// value is read last. If, for example, person was the first column listed
	// and it was null, the other columns would be ignored and Go would give
	// them zero values.
	rows, err := db.Query(`
		SELECT words.word, num, tense_id, person
		FROM verb_forms
		JOIN words on words.id = verb_forms.word_id
		WHERE inf_id = $1`,
		infId,
	)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	tenses := []*TenseInflection{&TenseInflection{}, &TenseInflection{}}
	for rows.Next() {
		var form string
		var person Person
		var num, tense int

		rows.Scan(&form, &num, &tense, &person)
		tense--
		if num == 2 {
			// plural won't have a person
			tenses[tense].Plural = append(tenses[tense].Plural, form)
		} else {
			if person&First != 0 {
				tenses[tense].First = append(tenses[tense].First, form)
			}
			if person&Second != 0 {
				tenses[tense].Second = append(tenses[tense].Second, form)
			}
			if person&Third != 0 {
				tenses[tense].Third = append(tenses[tense].Third, form)
			}
		}
	}

	return &ConjugationTable{
		Infinitive: inf,
		Present:    tenses[0],
		Past:       tenses[1],
	}, nil
}

// GetInfinitive retrieves the word and id of the verb form's infinitive.
func (db *PsqlDB) GetInfinitive(verbForm string) (string, int, error) {
	var infinitive string
	var id int
	err := db.QueryRow(`
		SELECT words.word, inf_id
		FROM verb_forms
		JOIN words on words.id = verb_forms.inf_id
		WHERE word_id = (SELECT id FROM words WHERE word = $1)`,
		verbForm,
	).Scan(&infinitive, &id)

	if err != nil {
		return "", 0, err
	}
	return infinitive, id, nil
}
