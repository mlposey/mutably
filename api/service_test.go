package main_test

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"log"
	"mutably/api"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	service.Router.ServeHTTP(recorder, request)
	return recorder
}

func checkCode(t *testing.T, expected, actual int) {
	t.Helper()
	if expected != actual {
		t.Error("Expected response code ", expected, ", got ", actual)
	}
}

// APIv1 should return a 404 response code if a specific language
// is requested but does not exist.
func TestGetLanguage_v1_missing(t *testing.T) {
	clearTable(t, "languages")

	req, _ := http.NewRequest("GET", "/api/v1/languages/3", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusNotFound, resp.Code)
}

// APIv1 should return a 200 response code along with the requested language
// object if it exists.
func TestGetLanguage_v1_exists(t *testing.T) {
	langId, _ := addLanguage(t)
	addLanguage(t)
	addLanguage(t)

	req, _ := http.NewRequest("GET", "/api/v1/languages/"+strconv.Itoa(langId),
		nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusOK, resp.Code)

	var respBody main.Language
	json.Unmarshal(resp.Body.Bytes(), &respBody)

	if respBody.Id != langId {
		t.Errorf("Expected id %d, got %d", langId, respBody.Id)
	}
}

// APIv1 should return a 404 response code if a collection of all languages
// is requested but the database has none.
func TestGetLanguages_v1_empty(t *testing.T) {
	clearTable(t, "languages")

	req, _ := http.NewRequest("GET", "/api/v1/languages", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusNotFound, resp.Code)
}

// APIv1 should return a 200 response code along with an array of all stored
// languages if any exist.
func TestGetLanguages_v1_notempty(t *testing.T) {
	clearTable(t, "languages")
	lang1, _ := addLanguage(t)
	lang2, _ := addLanguage(t)

	req, _ := http.NewRequest("GET", "/api/v1/languages", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusOK, resp.Code)

	var respBody []main.Language
	json.Unmarshal(resp.Body.Bytes(), &respBody)

	// Make sure it returns a correct array.
	lang1Count, lang2Count := 0, 0
	total := 0
	for _, language := range respBody {
		total++
		if language.Id == lang1 {
			lang1Count++
		}
		if language.Id == lang2 {
			lang2Count++
		}
	}

	if total != 2 {
		t.Errorf("Expected 2 results, got %d", total)
	}
	if lang1Count > 1 || lang2Count > 1 {
		t.Error("Seeing duplicate results for items that exist once")
	}
	if lang1Count < 1 || lang2Count < 1 {
		t.Error("Some languages weren't retrieved")
	}
}

// APIv1 should return a 404 response code if the database contains no words.
func TestGetWords_v1_empty(t *testing.T) {
	clearTable(t, "words")

	req, _ := http.NewRequest("GET", "/api/v1/words", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusNotFound, resp.Code)
}

// APIv1 should return a 200 response code along with an array of all words
// if any exist.
func TestGetWords_v1_notempty(t *testing.T) {
	clearTable(t, "words")
	wordId, _ := addWord(t)

	req, _ := http.NewRequest("GET", "/api/v1/words", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusOK, resp.Code)

	var words []main.Word
	json.Unmarshal(resp.Body.Bytes(), &words)

	if len(words) == 0 {
		t.Fatalf("Failed to retrieve words")
	}

	if words[0].Id != wordId {
		t.Errorf("Expected word id %d, got %d", wordId, words[0].Id)
	}
}

// APIv1 should return a 404 response code if a requested word does not exist.
func TestGetWord_v1_empty(t *testing.T) {
	clearTable(t, "words")

	req, _ := http.NewRequest("GET", "/api/v1/words/3", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusNotFound, resp.Code)
}

// APIv1 should return a 200 response code along with the requested word if
// it exists.
func TestGetWord_v1_notempty(t *testing.T) {
	clearTable(t, "words")
	wordId, _ := addWord(t)
	addWord(t)
	addWord(t)

	req, _ := http.NewRequest("GET", "/api/v1/words/"+strconv.Itoa(wordId), nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusOK, resp.Code)

	var word main.Word
	json.Unmarshal(resp.Body.Bytes(), &word)

	if word.Id != wordId {
		t.Errorf("Expected word id %d, got %d", wordId, word.Id)
	}
}

// APIv1 should return a 404 response code if the database contains no users.
func TestGetUsers_v1_empty(t *testing.T) {
	clearTable(t, "users")

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusNotFound, resp.Code)
}

// APIv1 should return a 200 response code along with an array of all users
// if any exist.
func TestGetUsers_v1_notempty(t *testing.T) {
	clearTable(t, "users")
	userId, _ := addUser(t)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusOK, resp.Code)

	var users []main.User
	json.Unmarshal(resp.Body.Bytes(), &users)

	if len(users) == 0 {
		t.Fatalf("Failed to retrieve users")
	}

	if users[0].Id != userId {
		t.Errorf("Expected user id %s, got %s", userId, users[0].Id)
	}
}

// APIv1 should return a 404 response code if a requested user does not exist.
func TestGetUser_v1_empty(t *testing.T) {
	clearTable(t, "users")

	req, _ := http.NewRequest("GET", "/api/v1/users/abd", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusNotFound, resp.Code)
}

// APIv1 should return a 200 response code along with the requested user if
// they exist.
func TestGetUser_v1_notempty(t *testing.T) {
	clearTable(t, "users")
	userId, _ := addUser(t)
	addUser(t)
	addUser(t)

	req, _ := http.NewRequest("GET", "/api/v1/users/"+userId, nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusOK, resp.Code)

	var user main.User
	json.Unmarshal(resp.Body.Bytes(), &user)

	if user.Id != userId {
		t.Errorf("Expected user id %s, got %s", userId, user.Id)
	}
}

// APIv1 should return a status created code if asked to create a user with
// name that is not in the database.
func TestCreateUser_v1_unique(t *testing.T) {
	clearTable(t, "users")

	req, _ := http.NewRequest("POST", "/api/v1/users", nil)
	cred := base64.StdEncoding.EncodeToString([]byte("user:pass"))
	req.Header.Set("Authorization", "Basic "+cred)

	resp := sendRequest(req)
	checkCode(t, http.StatusCreated, resp.Code)
}

// APIv1 should return a bad request code if asked to create a user with a
// name that already exists.
func TestCreateUser_v1_duplicate(t *testing.T) {
	clearTable(t, "users")

	var resp *httptest.ResponseRecorder
	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest("POST", "/api/v1/users", nil)
		cred := base64.StdEncoding.EncodeToString([]byte("user:pass"))
		req.Header.Set("Authorization", "Basic "+cred)

		resp = sendRequest(req)
	}
	checkCode(t, http.StatusBadRequest, resp.Code)
}

// TODO: Test POST /users bad or missing authorization header.
// TODO: Test POST /users bad username:password format

func clearTable(t *testing.T, table string) {
	t.Helper()
	_, err := db.Exec("DELETE FROM " + table)
	if err != nil {
		t.Error(err)
	}
}

func addUser(t *testing.T) (string, string) {
	t.Helper()

	userPassword := uuid.NewV4().String()
	var userId string

	err := db.QueryRow(`
		SELECT create_user($1, $2)`,
		uuid.NewV4().String(), userPassword,
	).Scan(&userId)
	if err != nil {
		t.Error(err)
	}
	return userId, userPassword
}

func addWord(t *testing.T) (int, int) {
	t.Helper()

	langId, _ := addLanguage(t)
	word := uuid.NewV4().String()

	var tableId int
	err := db.QueryRow(`
		SELECT add_infinitive($1, $2)`,
		word, langId,
	).Scan(&tableId)
	if err != nil {
		t.Error(err)
	}

	var wordId int
	err = db.QueryRow(`
		SELECT word_id FROM verbs
		WHERE  conjugation_table = $1`,
		tableId,
	).Scan(&wordId)
	if err != nil {
		t.Error(err)
	}

	return wordId, langId
}

func addLanguage(t *testing.T) (id int, name string) {
	t.Helper()
	name = uuid.NewV4().String()

	err := db.QueryRow(`
		INSERT INTO languages (language, tag)
		VALUES ($1, $2)
		RETURNING id`,
		name, uuid.NewV4().String(),
	).Scan(&id)

	if err != nil {
		t.Error(err)
	}
	return id, name
}
