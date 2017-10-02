// helper and setup functions for the tests
package main_test

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"mutably/api"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var service *main.Service
var db *sql.DB

func init() {
	host := os.Getenv("DATABASE_HOST")
	name := os.Getenv("DATABASE_NAME")
	user := os.Getenv("DATABASE_USER")
	pwd := os.Getenv("DATABASE_PASSWORD")

	var err error
	db, err = sql.Open("postgres", fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s sslmode=disable",
		name, user, pwd, host,
	))

	// This won't directly be used by tests.
	database, err := main.NewDB(host, name, user, pwd)
	if err != nil {
		log.Fatal("Could not access database; ", err)
	}

	service, err = main.NewService(database, "8080")
	if err != nil {
		log.Fatal("Could not start service; ", err)
	}
}

// sendRequest submits and then records the result of an HTTP request.
func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	service.Router.ServeHTTP(recorder, request)
	return recorder
}

func checkError(t *testing.T, e error) {
	if e != nil {
		t.Error(e)
	}
}

// checkCode performs an error check on two response codes that should be equal.
func checkCode(t *testing.T, expected, actual int) {
	t.Helper()
	if expected != actual {
		t.Error("Expected response code ", expected, ", got ", actual)
	}
}

// clearTable removes all items from a specific table in the test database.
func clearTable(t *testing.T, table string) {
	t.Helper()
	_, err := db.Exec("DELETE FROM " + table)
	checkError(t, err)
}

// clearDatabase returns the test database to the default empty state.
func clearDatabase(t *testing.T) {
	t.Helper()
	clearTable(t, "verb_forms")
	clearTable(t, "words")
	clearTable(t, "languages")
	clearTable(t, "users")
}

// createCompleteVerb inserts the complete structure of a made up verb.
// That structure includes the infinitive and number/singularity of present
// and past tenses.
// returns (infinitive, id of infinitive)
func createCompleteVerb(t *testing.T) (string, int) {
	t.Helper()
	var langId int
	db.QueryRow(`
		INSERT INTO languages (name)
		VALUES ('dutch') RETURNING id`,
	).Scan(&langId)

	var infId int
	db.QueryRow(`INSERT INTO words (word) VALUES ('krijgen') RETURNING id`).Scan(&infId)

	db.Exec(`
		INSERT INTO words (word) VALUES
		('krijg'), ('krijgt'), ('kreeg'), ('kreegt'), ('kregen')
	`)
	db.Exec(`
		INSERT INTO verb_forms (lang_id, word_id, inf_id, tense_id, person, num)
		VALUES
		($1, (SELECT id FROM words WHERE word = 'krijg'), $2, 1, 2, 1),
		($1, (SELECT id FROM words WHERE word = 'krijgt'), $2, 1, 12, 1),
		($1, $2, $2, 1, NULL, 2),
		($1, (SELECT id FROM words WHERE word = 'kreeg'), $2, 2, 14, 1),
		($1, (SELECT id FROM words WHERE word = 'kreegt'), $2, 2, 4, 1),
		($1, (SELECT id FROM words WHERE word = 'kregen'), $2, 2, NULL, 2)`,
		langId, infId,
	)
	return "krijgen", infId
}

// createVerbForm inserts a value into the test database's verb_forms table.
// returns (language id, word id, verb form id)
func createVerbForm(t *testing.T) (int, int, int) {
	t.Helper()

	_, wordId := addWord(t)
	langId, _ := addLanguage(t)

	var formId int
	err := db.QueryRow(`
		INSERT INTO verb_forms (lang_id, word_id, inf_id, tense_id, num)
		VALUES ($1, $2, $3, (SELECT id FROM tenses WHERE tense = 'past'), 2)
		RETURNING id`,
		langId, wordId, wordId,
	).Scan(&formId)

	checkError(t, err)
	return langId, wordId, formId
}

// addWord inserts a word into the test database's words table.
// returns (the inserted word, word's id)
func addWord(t *testing.T) (string, int) {
	t.Helper()

	word := uuid.NewV4().String()
	var wordId int

	err := db.QueryRow(`
		INSERT INTO words (word)
		VALUES ($1) RETURNING id`,
		word,
	).Scan(&wordId)

	checkError(t, err)
	return word, wordId
}

// addLanguage inserts a language into the test database's languages table.
// returns (the language's id, the language)
func addLanguage(t *testing.T) (id int, name string) {
	t.Helper()
	name = uuid.NewV4().String()

	err := db.QueryRow(`
		INSERT INTO languages (name, tag)
		VALUES ($1, $2)
		RETURNING id`,
		name, uuid.NewV4().String(),
	).Scan(&id)

	checkError(t, err)
	return id, name
}

// requestAccount supplies credentials to the HTTP account creation resource.
// Returns the user's jwt token
func requestAccount(t *testing.T, user, pass string) string {
	t.Helper()

	req, err := http.NewRequest("POST", "/api/v1/users", nil)
	checkError(t, err)

	cred := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	req.Header.Set("Authorization", "Basic "+cred)
	resp := sendRequest(req)
	checkCode(t, http.StatusCreated, resp.Code)

	var respBody map[string]string
	json.Unmarshal(resp.Body.Bytes(), &respBody)

	return respBody["token"]
}

// makeAdmin gives a user the admin role.
func makeAdmin(t *testing.T, userId string) {
	t.Helper()
	_, err := db.Exec(`
		UPDATE users SET role_id = (
			SELECT id FROM roles WHERE role = 'admin'
		) WHERE id = $1`,
		userId,
	)
	checkError(t, err)
}

// createUser inserts a new user into the test database.
// returns (user id, username, password)
func createUser(t *testing.T) (string, string, string) {
	t.Helper()
	username := uuid.NewV4().String()
	password := uuid.NewV4().String()
	var userId string

	err := db.QueryRow(`
		SELECT create_user($1, $2)`,
		username, password,
	).Scan(&userId)

	checkError(t, err)
	return userId, username, password
}
