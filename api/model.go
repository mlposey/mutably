package main

import (
	"database/sql"
	"time"
)

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
	Tag sql.NullString `json:"tag"`
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

// CreateUser inserts a new user into the database.
// Duplicate names are not allowed; attempting to insert one will result
// in a non-nil error.
func (db *PsqlDB) CreateUser(name, password string) (string, error) {
	var userId string
	err := db.QueryRow(`SELECT create_user($1, $2)`, name, password).Scan(&userId)
	return userId, err
}
