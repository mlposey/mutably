package controller

import (
	"mutably/api/model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Words is a Controller that handles the /words resource.
type Words struct {
	db model.Database
}

func (w *Words) Routes() []Route {
	return []Route{
		{ // GET /v1/words
			Version:     "v1",
			Path:        "/words",
			Method:      "GET",
			Handler:     w.getWords,
			IsProtected: false,
		},
		{ // GET /v1/words/{id:[0-9]+}
			Version:     "v1",
			Path:        "/words/{id:[0-9]+}",
			Method:      "GET",
			Handler:     w.getWord,
			IsProtected: false,
		},
		{ // GET /v1/words/{word}/inflections
			Version:     "v1",
			Path:        "/words/{word}/inflections",
			Method:      "GET",
			Handler:     w.getInflections,
			IsProtected: false,
		},
	}
}

// GET /api/v1/words
func (ws *Words) getWords(w http.ResponseWriter, r *http.Request) {
	words, err := ws.db.GetWords()
	respondWithAggregate(w, words, len(words), err, "words")
}

// GET /api/v1/words/{id:[0-9]+}
func (ws *Words) getWord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wordId, err := strconv.Atoi(vars["id"])
	if err != nil {
		makeErrorResponse(w, http.StatusBadRequest, "invalid word id")
	}

	word, err := ws.db.GetWord(wordId)
	if err != nil {
		makeErrorResponse(w, http.StatusNotFound, "word not found")
	} else {
		makeJsonResponse(w, http.StatusOK, word)
	}
}

// GET /api/v1/words/{word}/inflections
func (ws *Words) getInflections(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table, err := ws.db.GetConjugationTable(vars["word"])
	if err != nil {
		makeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	makeJsonResponse(w, http.StatusOK, table)
}
