package main_test

import (
	"database/sql"
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

func clearLanguages(t *testing.T) {
	t.Helper()
	_, err := db.Exec("DELETE FROM languages")
	if err != nil {
		t.Error(err)
	}
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

// APIv1 should return a 404 response code if a specific language
// is requested but does not exist.
func TestGetLanguage_v1_missing(t *testing.T) {
	clearLanguages(t)

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

	var respBody map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &respBody)

	if respBody["id"] != langId {
		t.Errorf("Expected id %d, got %d", langId, respBody["id"])
	}
}

// APIv1 should return a 404 response code if a collection of all languages
// is requested but the database has none.
func TestGetLanguages_v1_empty(t *testing.T) {
	clearLanguages(t)

	req, _ := http.NewRequest("GET", "/api/v1/languages", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusNotFound, resp.Code)
}

// APIv1 should return a 200 response code along with an array of all stored
// languages if any exist.
func TestGetLanguages_v1_notempty(t *testing.T) {
	clearLanguages(t)
	lang1, _ := addLanguage(t)
	lang2, _ := addLanguage(t)

	req, _ := http.NewRequest("GET", "/api/v1/languages", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusOK, resp.Code)

	var respBody []map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &respBody)

	// Make sure it returns a correct array.
	lang1Count, lang2Count := 0, 0
	total := 0
	for i := range respBody {
		total++
		if respBody[i]["id"] == lang1 {
			lang1Count += 1
		}
		if respBody[i]["id"] == lang2 {
			lang2Count += 1
		}
	}

	if total != 2 {
		t.Errorf("Expected 2 results, got %d", total)
	}
	if lang1Count != 1 || lang2Count != 1 {
		t.Errorf("Seeing duplicate results for items that exist once")
	}
}
