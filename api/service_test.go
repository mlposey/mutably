package main_test

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"mutably/api"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// APIv1 should return a 404 response code if a specific language
// is requested but does not exist.
func TestGetLanguage_v1_missing(t *testing.T) {
	clearDatabase(t)

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
	clearDatabase(t)

	req, _ := http.NewRequest("GET", "/api/v1/languages", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusNotFound, resp.Code)
}

// APIv1 should return a 200 response code along with an array of all stored
// languages if any exist.
func TestGetLanguages_v1_notempty(t *testing.T) {
	clearDatabase(t)
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
	clearDatabase(t)

	req, _ := http.NewRequest("GET", "/api/v1/words", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusNotFound, resp.Code)
}

// APIv1 should return a 200 response code along with an array of all words
// if any exist.
func TestGetWords_v1_notempty(t *testing.T) {
	clearDatabase(t)
	_, wordId, _ := createVerbForm(t)

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
	clearDatabase(t)

	req, _ := http.NewRequest("GET", "/api/v1/words/3", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusNotFound, resp.Code)
}

// APIv1 should return a 200 response code along with the requested word if
// it exists.
func TestGetWord_v1_notempty(t *testing.T) {
	clearDatabase(t)
	_, wordId, _ := createVerbForm(t)
	createVerbForm(t)
	createVerbForm(t)

	req, _ := http.NewRequest("GET", "/api/v1/words/"+strconv.Itoa(wordId), nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusOK, resp.Code)

	var word main.Word
	json.Unmarshal(resp.Body.Bytes(), &word)

	if word.Id != wordId {
		t.Errorf("Expected word id %d, got %d", wordId, word.Id)
	}
}

// APIv1 should return a 401 response code if the client sends a request
// and has no token.
func TestGetUsers_v1_forbidden(t *testing.T) {
	clearDatabase(t)
	createUser(t)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusUnauthorized, resp.Code)
}

// APIv1 should return a 401 response code if the client sends a request
// and has no token.
func TestGetUser_v1_forbidden(t *testing.T) {
	clearDatabase(t)
	userId, _, _ := createUser(t)

	req, _ := http.NewRequest("GET", "/api/v1/users/"+userId, nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusUnauthorized, resp.Code)
}

// APIv1 should return a status created code if asked to create a user with
// name that is not in the database.
func TestCreateUser_v1_unique(t *testing.T) {
	clearDatabase(t)
	requestAccount(t, "user", "pass")
}

// APIv1 should return a bad request code if asked to create a user with a
// name that already exists.
func TestCreateUser_v1_duplicate(t *testing.T) {
	clearDatabase(t)

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

// APIv1 should a not found request code if asked to retrieve the inflection
// table of a verb that does not exist.
func TestGetInflections_v1_missing(t *testing.T) {
	clearDatabase(t)

	word := uuid.NewV4().String()
	req, _ := http.NewRequest("GET", "/api/v1/words/"+word+"/inflections", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusNotFound, resp.Code)
}

// APIv1 should return the inflection table of verbs that exist.
// The table is expected to look like the following:
// {
//	"infinitive": string,
//  "present": {"first":[string...], "second":[string...], "third":[string...], "plural":[string...]},
//  "past": {"first":[string...], "second":[string...], "third":[string...], "plural":[string...]}
// }
func TestGetInflections_v1_exists(t *testing.T) {
	clearDatabase(t)
	// We'll request a verb that has forms in each grammatical category. Thus,
	// the resulting JSON response should have entries in each position.
	infinitive, _ := createCompleteVerb(t)

	req, _ := http.NewRequest("GET", "/api/v1/words/"+infinitive+"/inflections", nil)
	resp := sendRequest(req)
	checkCode(t, http.StatusOK, resp.Code)

	var table main.ConjugationTable
	json.Unmarshal(resp.Body.Bytes(), &table)

	missingDataError := errors.New("Table is missing category information")
	if table.Infinitive == "" || table.Present == nil || table.Past == nil {
		t.Error(missingDataError)
	}

	tenses := []*main.TenseInflection{table.Present, table.Past}
	for _, tense := range tenses {
		if len(tense.First) == 0 ||
			len(tense.Second) == 0 ||
			len(tense.Third) == 0 ||
			len(tense.Plural) == 0 {
			t.Error(missingDataError)
		}
	}
}

// TODO: Test GET /words/{word}/inflections using all forms of a verb.
//       Right now the tests only check the infinitive, but we should ensure
//       that API calls that use the various forms also retrieve the same
//       table.

// APIv1 should return a new JWT if /tokens is given a set of credentials
// (through the Authorization header) that match an existing user.
func TestGetSession_validCredentials(t *testing.T) {
	clearDatabase(t)

	req, err := http.NewRequest("GET", "/api/v1/tokens", nil)
	checkError(t, err)

	_, user, pass := createUser(t)
	cred := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	req.Header.Set("Authorization", "Basic "+cred)

	resp := sendRequest(req)
	if resp.Code != http.StatusOK {
		t.Error("Expected 200 status when providing valid credentials")
	}

	var respBody map[string]string
	json.Unmarshal(resp.Body.Bytes(), &respBody)
	if respBody["token"] == "" {
		t.Error("Expected token in body of validated /tokens response")
	}
}

// APIv1 should return an error message if /tokens is given a set of credentials
// that do not belong to any user.
func TestGetSession_invalidCredentials(t *testing.T) {
	clearDatabase(t)

	req, err := http.NewRequest("GET", "/api/v1/tokens", nil)
	checkError(t, err)

	user, pass := "not_a_real_user", "a_fake_pass"
	cred := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	req.Header.Set("Authorization", "Basic "+cred)

	resp := sendRequest(req)
	if resp.Code == http.StatusOK {
		t.Error("Found 200 status code in result of failed GET /tokens")
	}

	var respBody map[string]string
	json.Unmarshal(resp.Body.Bytes(), &respBody)
	if respBody["error"] == "" {
		t.Error("Expected error response from bad GET /tokens")
	}
}